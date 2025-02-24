package apiserver

type Config struct {
	BindPort string `yaml:"port" env:"HELLO_PORT" default:"9001"`
	BindHost string `yaml:"host" env:"HELLO_HOST" default:"0.0.0.0"`
}

func NewConfig() *Config {
	return &Config{
		BindPort: "",
		BindHost: "0.0.0.0", // костылек, но host нигде не указывается в задании
	}
}

func DefaultConfig() *Config {
	return &Config{
		BindPort: "9001",
		BindHost: "0.0.0.0",
	}
}
