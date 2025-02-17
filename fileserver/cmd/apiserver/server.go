package main

import (
	"flag"
	"log"

	"yadro.com/course/internal/apiserver"
	"yadro.com/course/internal/storage"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultConfigPath  = "config.yaml"
	defaultStoragePath = "./data"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to config file")
	flag.Parse()
}

func getConfig() *apiserver.Config {
	config := apiserver.NewConfig()

	if err := cleanenv.ReadEnv(config); err == nil {
		log.Println("Using environment variables for application")
		return config
	}

	log.Printf("Cannot find environment variables, trying config file %s....", configPath)

	if err := cleanenv.ReadConfig(configPath, config); err == nil {
		log.Printf("Using config file with path %s", config)
		return config
	}

	log.Println("failed to load configuration from file or environment variables")

	return apiserver.NewConfig()
}

func main() {
	config := getConfig()
	fStorage, err := storage.NewStorage(defaultStoragePath)
	if err != nil {
		log.Panicf("Error while creating storage with path %s: %v", defaultStoragePath, err)
	}

	s := apiserver.NewServer(config, fStorage)
	s.Run()
}
