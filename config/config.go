package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Backend struct {
	URL string `yaml:"url"`
}

type Config struct {
	Port     string    `yaml:"port"`
	Strategy string    `yaml:"strategy"`
	Backends []Backend `yaml:"backends"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
