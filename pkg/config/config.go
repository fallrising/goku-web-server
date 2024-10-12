package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port   int    `mapstructure:"PORT"`
	DBPath string `mapstructure:"DB_PATH"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
