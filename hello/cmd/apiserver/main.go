package main

import (
	"flag"
	"fmt"
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

func parseConfig(configPath string) (*apiserver.Config, error) {
	config := apiserver.NewConfig()

	if err := cleanenv.ReadConfig(configPath, config); err == nil {
		log.Printf("Using config file: %s", configPath)
		return config, nil
	}
	log.Printf("Cannot read config file %s, trying environment variables...", configPath)

	if err := cleanenv.ReadEnv(config); err == nil {
		log.Println("Using environment variables for configuration")
		return config, nil
	}

	return nil, fmt.Errorf("failed to load configuration from file or environment variables, using default")
}

func main() {
	config, err := parseConfig(configPath)
	if err != nil {
		log.Panicf("Error while reading a config: %v", err)
	}

	s := apiserver.NewServer(config)
	s.Run()
}
