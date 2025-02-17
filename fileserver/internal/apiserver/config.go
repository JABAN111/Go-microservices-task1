package apiserver

type Config struct {
	BindPort   string `yaml:"port" env:"FILESERVER_PORT" default:"9001"`
	BindHost   string `yaml:"host" env:"FILESERVER_HOST" default:"0.0.0.0"`
	ConfigPath string `yaml:"path" env:"FILESERVER_CONFIG_PATH" default:"./data"`
}

func NewConfig() *Config {
	return &Config{
		BindPort:   "",
		BindHost:   "0.0.0.0", //костылек, в задании не задается хост
		ConfigPath: "./data",  //костылек, в задании не задается путь
	}
}

func DefaultConfig() *Config {
	return &Config{
		BindPort:   "9001",
		BindHost:   "0.0.0.0",
		ConfigPath: "./data",
	}
}
