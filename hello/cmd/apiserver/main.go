package main

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
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
}

func parseConfig() (*apiserver.Config, error) {
	config := apiserver.NewConfig()

	if err := cleanenv.ReadConfig(configPath, config); err == nil {
		log.Printf("Using config file with path %s", config)
		return config, nil
	}

	log.Printf("Cannot read config file %s, trying environment variables...", configPath)

	if err := cleanenv.ReadEnv(config); err == nil {
		log.Println("Using environment variables for application")
		return config, nil
	}
	err := fmt.Errorf("failed to load configuration from file or environment variables")
	log.Println(err)

	return nil, err
}

func main() {
	flag.Parse()

	config, err := parseConfig()
	if err != nil {
		log.Panicf("Error while reading a config: %v", err)
	}

	s := apiserver.NewServer(config)
	s.Run()
}
