package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/BurntSushi/toml"
)

// Config represents database server and credentials
type Config struct {
	JWTSecret     string
	MongoServer   string
	MongoDatabase string
	PublicDir     string
}

var instance *Config
var once sync.Once

// GetInstance return singleton of config
func GetInstance() *Config {
	once.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		configPath := fmt.Sprintf("%s/../../settings/config.toml", filepath.Dir(filename))
		var conf Config
		if _, err := toml.DecodeFile(configPath, &conf); err != nil {
			log.Fatal(err)
		}
		instance = &conf
	})
	return instance
}
