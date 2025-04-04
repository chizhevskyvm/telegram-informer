package config

import (
	"github.com/go-yaml/yaml"
	"os"
)

type Config struct {
	Telegram       Telegram       `yaml:"telegram"`
	DatabaseConfig DatabaseConfig `yaml:"db"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

type DatabaseConfig struct {
	Connection string `yaml:"connection"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
