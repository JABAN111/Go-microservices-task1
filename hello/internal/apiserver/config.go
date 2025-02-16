package apiserver

type Config struct {
	BindPort string `yaml:"port" env:"HELLO_PORT" default:"9001"`
	BindHost string `yaml:"host" env:"HELLO_HOST" default:"localhost"`
}

func NewConfig() *Config {
	return &Config{
		BindPort: "9001",
		BindHost: "0.0.0.0",
	}
}
