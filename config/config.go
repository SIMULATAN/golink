package config

import (
	"golink/service/link"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Postgres link.PostgresConfig `yaml:"postgres"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	io, err := os.ReadFile("config.yml")
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(io, &cfg)
	return cfg, err
}
