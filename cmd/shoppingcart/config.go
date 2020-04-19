package main

import (
	"github.com/kelseyhightower/envconfig"
)

// DatabaseConfig defines a configuration for database connection.
type DatabaseConfig struct {
	User     string `envconfig:"database_user"`
	Password string `envconfig:"database_password"`
	Database string `envconfig:"database_name"`
	Host     string `envconfig:"database_host"`
	Port     int    `envconfig:"database_port"`
}

// Config describes the relevant settings from environment variables.
type Config struct {
	Database DatabaseConfig
}

// NewConfig returns a Config which is populated by environment variables.
func NewConfig() (*Config, error) {
	config := &Config{}

	if err := envconfig.Process("shoppingcart", config); err != nil {
		return nil, err
	}

	return config, nil
}
