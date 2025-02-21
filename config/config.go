package config

import (
	"github.com/spf13/viper"
	"golink/service/link"
	"strings"
)

type Config struct {
	BaseUrl  string              `yaml:"base_url" mapstructure:"base_url"`
	Postgres link.PostgresConfig `yaml:"postgres" mapstructure:"postgres"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("golink")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
