package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string       `yaml:"env" env-required:"true"`
	DbConnection dbConnection `yaml:"db_connection" env-required:"true"`
	HttpServer   HttpServer   `yaml:"http_server" env-required:"true"`
}

type dbConnection struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Dbname   string `yaml:"dbname" env-required:"true"`
}

type HttpServer struct {
	Address      string        `yaml:"address" env-default:"localhost:8081"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	IddleTimeout string        `yaml:"iddle_timeout" env-default:"60s"`
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
