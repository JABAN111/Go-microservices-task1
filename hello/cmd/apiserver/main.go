package main

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"yadro.com/course/internal/apiserver"
)

const (
	defaultConfigPath = "config.yaml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to config file")
	flag.Parse()
}

func getConfig(configPath string) *apiserver.Config {
	config := apiserver.NewConfig()

	if err := cleanenv.ReadEnv(config); err == nil {
		log.Println("Using environment variables for configuration")
		return config
	}

	log.Printf("Cannot find environment variables, trying config file %s...", configPath)

	if err := cleanenv.ReadConfig(configPath, config); err == nil {
		log.Printf("Using config file: %s", configPath)
		return config
	}

	log.Printf("Failed to load configuration from file or environment variables, using default configuration")
	return apiserver.NewConfig()
}

func main() {
	config := getConfig(configPath)
	s := apiserver.NewServer(config)
	s.Run()
}
