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

	if err := cleanenv.ReadConfig(configPath, config); err == nil {
		if config.BindPort != "" {
			log.Printf("Using configuration: %s", config)
			return config
		}
	}

	log.Printf("Failed to load configuration from file or environment variables, using default configuration")
	return apiserver.DefaultConfig()
}

func main() {
	config := getConfig(configPath)
	s := apiserver.NewServer(config)
	s.Run()
}
