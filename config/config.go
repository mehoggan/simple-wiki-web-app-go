package config

import (
	"log"
	"os"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
)

func loadConfig(resourcePath string) *types.Config {
	file, err := os.Open(resourcePath)
	if err != nil {
		log.Fatalf("Failed to open %s!!!", resourcePath)
	} else {
		log.Printf("Opened %s for reading, loading...", resourcePath)
	}

	var config types.Config
	log.Printf("Decoding settings...")
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config from %s with %s!!!", resourcePath, err)
	}
	return &config
}

var instantiated *types.Config
var once sync.Once

func Intantiate(configPath string) *types.Config {
	once.Do(func() {
		instantiated = loadConfig(configPath)
	})
	return instantiated
}
