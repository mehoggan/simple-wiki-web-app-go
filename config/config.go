package config

import (
	"log"
	"os"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

func loadConfig(resourcePath string) *Config {
	file, err := os.Open(resourcePath)
	if err != nil {
		log.Fatalf("Failed to open (%v)!!!", resourcePath)
	} else {
		log.Printf("Opened %v for reading, loading...", resourcePath)
	}

	var config Config
	log.Printf("Decoding settings...")
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to load config from (%v)!!!", resourcePath)
	} else {
		log.Printf("Decoding complete.")
	}
	return &config
}

var instantiated *Config
var once sync.Once

func Intantiate(configPath string) *Config {
	once.Do(func() {
		instantiated = loadConfig(configPath)
	})
	return instantiated
}
