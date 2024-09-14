package database

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB ConfigDatabase `yaml:"db"`
}

type ConfigDatabase struct {
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"USER" env-default:"user"`
	Password string `yaml:"password" env:"PASSWORD"`
}

func Load(cfgPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}
	return &cfg
}
