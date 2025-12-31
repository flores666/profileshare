package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string     `yaml:"env" env-required:"true"`
	ConnectionString string     `yaml:"db_connection_string" env-required:"true"`
	HttpServer       HttpServer `yaml:"http_server" env-required:"true"`
	Producer         Producer   `yaml:"producer" env-required:"true"`
}

type HttpServer struct {
	Address      string        `yaml:"address" env-default:"localhost:8081"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type Producer struct {
	Brokers []string `yaml:"brokers" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalln("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalln("config file does not exist: ", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalln("error when reading config: ", err)
	}

	return cfg
}
