package server

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App ConfigApp `yaml:"app"`
}

type ConfigApp struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func Load(cfgPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}
	return &cfg
}
