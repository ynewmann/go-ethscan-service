package useful_servise

import (
	"bytes"
	"encoding/json"
	"go-ethscan-service/etherscan"
	"log"
	"os"
	"path"
)

const DefaultRelativePath = ".jmind"
const DefaultConfigPath = ".jmind/config.cfg"
const DefaultMemoryCacheBackupPath = ".jmind/cache.backup"

func DefaultConfig() *Config {
	return &Config{
		UseCache: false,
		ApiUrl:   etherscan.DefaultApiUrl,
	}
}

type Config struct {
	UseCache              bool
	MemoryCacheBackupPath string
	CacheSize             uint32

	Address string

	ApiUrl string
	ApiKey string
}

func (cfg *Config) SaveToFile(p string) error {
	_, err := os.Stat(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		dir := path.Dir(p)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&cfg)
	if err != nil {
		return err
	}

	_, err = f.Write(b.Bytes())
	return err
}

func (cfg *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	b := bytes.NewReader(data)
	return json.NewDecoder(b).Decode(&cfg)
}
