package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Detect map[string]Target `yaml:"detect"`
}

type Target struct {
	Mode       string   `yaml:"mode"`
	Enviroment string   `yaml:"environment"`
	Paths      []string `yaml:"paths"`
}

func FromFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	cfg := &Config{}
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config yaml: %w", err)
	}

	return cfg, nil
}
