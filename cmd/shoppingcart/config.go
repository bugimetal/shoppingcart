package main

import (
	"github.com/kelseyhightower/envconfig"
)

type DatabaseConfig struct {
	User     string `envconfig:"database_user"`
	Password string `envconfig:"database_password"`
	Database string `envconfig:"database_name"`
	Host     string `envconfig:"database_host"`
	Port     int    `envconfig:"database_port"`
}

type Config struct {
	Database DatabaseConfig
}

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := envconfig.Process("shoppingcart", config); err != nil {
		return nil, err
	}

	return config, nil
}
