package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HttpServer `yaml:"http_server"`
	DB         `yaml:"db"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"360s"`
}

type DB struct {
	Host           string        `yaml:"host" env-default:"localhost"`
	Port           int           `yaml:"port" env-default:"5432"`
	User           string        `yaml:"user" env-default:"postgres"`
	Password       string        `yaml:"password" env-default:"postgres"`
	Name           string        `yaml:"name" env-default:"postgres"`
	Sslmode        string        `yaml:"sslmode" env-default:"disable"`
	MaxAttempts    int           `yaml:"max_attempts" env-default:"5"`
	Delay          time.Duration `yaml:"delay" env-default:"5s"`
	MigrationsPath string        `yaml:"migrations_path" env-default:"migrations"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalln("CONFIG_PATH does not exist")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalln("Can not read config file")
	}

	return &cfg
}
