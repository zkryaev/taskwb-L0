package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB  ConfigDB  `yaml:"db"`
	App ConfigApp `yaml:"app"`
}

type ConfigApp struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ConfigDB struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func Load(cfgPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}
	return &cfg
}
