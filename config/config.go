package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Telegram Telegram       `yaml:"telegram"`
	Database DatabaseConfig `yaml:"db"`
	Redis    RedisConfig    `yaml:"redis"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

type DatabaseConfig struct {
	Connection string `yaml:"connection"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
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
