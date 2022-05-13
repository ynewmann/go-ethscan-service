package useful_servise

import (
	"log"
	"os"
	"os/signal"

	"go-ethscan-service/etherscan"
	"go-ethscan-service/storage"

	"github.com/valyala/fasthttp"
)

type (
	Option func(*UsefulService)

	UsefulService struct {
		api              etherscan.Api
		server           *fasthttp.Server
		totalAmountCache storage.Storage
		cfg              *Config
	}
)

func WithStorage(s storage.Storage) Option {
	return func(us *UsefulService) {
		us.totalAmountCache = s
	}
}

func WithApi(api etherscan.Api) Option {
	return func(us *UsefulService) {
		us.api = api
	}
}

func WithConfig(cfg *Config) Option {
	return func(us *UsefulService) {
		us.cfg = cfg
	}
}

func NewUsefulService(options ...Option) *UsefulService {
	us := &UsefulService{
		server: &fasthttp.Server{},
		cfg:    &Config{Address: ":3333"},
	}

	us.server.Handler = newMiddleware(newHandler(us))
	for _, opt := range options {
		opt(us)
	}

	if us.cfg.UseCache && us.totalAmountCache == nil {
		log.Println("Using default totalAmountCache...")
		us.totalAmountCache = storage.NewMemoryCache(us.cfg.CacheSize)
	}

	if us.cfg.UseCache && us.cfg.MemoryCacheBackupPath != "" {
		s, ok := us.totalAmountCache.(storage.MemoryStorage)
		if ok {
			log.Println("Trying to load totalAmountCache backup...")
			err := s.LoadFromFile(us.cfg.MemoryCacheBackupPath)
			if err != nil {
				log.Printf("Failed to load totalAmountCache backup: %s", err.Error())
			}
		}
	}

	if us.api == nil {
		opts := make([]etherscan.ApiClientOption, 0, 2)
		if us.cfg.ApiUrl != "" {
			opts = append(opts, etherscan.WithUrl(us.cfg.ApiUrl))
		}
		if us.cfg.ApiKey != "" {
			opts = append(opts, etherscan.WithApiKey(us.cfg.ApiKey))
		}

		us.api = etherscan.NewApi(etherscan.WithClient(etherscan.NewApiClient(opts...)))
	}

	return us
}

func (us *UsefulService) Start() error {
	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt, os.Kill)
		<-stop

		err := us.Shutdown()
		if err != nil {
			log.Println(err)
		}
		close(stop)
	}()

	log.Printf("Start listening at %s...", us.cfg.Address)
	return us.server.ListenAndServe(us.cfg.Address)
}

func (us *UsefulService) Shutdown() error {
	s, ok := us.totalAmountCache.(storage.MemoryStorage)
	if ok && us.cfg.UseCache && us.cfg.MemoryCacheBackupPath != "" {
		err := s.SaveToFile(us.cfg.MemoryCacheBackupPath)
		if err != nil {
			log.Printf("Failed to save memory storage: %s", err.Error())
		}
	}

	return us.server.Shutdown()
}
