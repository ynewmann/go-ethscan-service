package main

import (
	"flag"
	us "go-ethscan-service/useful-servise"
	"log"
	"os"
	"path/filepath"
)

func main() {
	apiUrl := flag.String("ApiUrl", "", "Api URL")
	apiKey := flag.String("ApiKey", "", "Api key")
	useCache := flag.Bool("useCache", true, "use cache")
	port := flag.String("port", ":3333", "server address")
	cacheSize := flag.Uint("cacheSize", 1000, "server address")
	memoryCacheBackupPath := flag.String("cache-bckp", "", "backup of cache")
	loadCfg := flag.String("load-cfg", "", "config path")

	flag.Parse()

	var cfg *us.Config
	if *loadCfg != "" {
		log.Fatalln(cfg.LoadFromFile(*loadCfg))
	} else {
		cfg = &us.Config{
			UseCache:              *useCache,
			Address:               *port,
			MemoryCacheBackupPath: *memoryCacheBackupPath,
			CacheSize:             uint32(*cacheSize),
			ApiUrl:                *apiUrl,
			ApiKey:                *apiKey,
		}

		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}

		if cfg.UseCache && cfg.MemoryCacheBackupPath == "" {
			cfg.MemoryCacheBackupPath = filepath.Join(dirname, us.DefaultMemoryCacheBackupPath)
		}

		path := filepath.Join(dirname, us.DefaultConfigPath)
		err = cfg.SaveToFile(path)
		if err != nil {
			log.Fatalln(err)
		}
	}

	server := us.NewUsefulService(us.WithConfig(cfg))
	log.Fatalln(server.Start())
}
