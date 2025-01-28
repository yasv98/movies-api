package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Port     string  `yaml:"port" validate:"required"`
		MonogoDB MongoDB `yaml:"mongodb" validate:"required"`
	}

	MongoDB struct {
		URI      string `yaml:"uri" validate:"required"`
		Database string `yaml:"database" validate:"required"`
	}
)

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}
