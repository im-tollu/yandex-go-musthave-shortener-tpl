// Package config is responsible for taking the runtime configuration from
// multiple sources of parameters and providing a structured configuration
// data to the service at the time of launch. It is also provides sensible
// defaults.
//
// Environment variables are considered the primary source of configuration.
// It supports the 12-factors app approach.
// For developers' convenience configuration can be overridden
// with CLI parameters.
package config

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string  `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       url.URL `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DatabaseDSN   string  `env:"DATABASE_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"`
}

func Load() (*Config, error) {
	config := &Config{}

	if errEnv := env.Parse(config); errEnv != nil {
		return nil, fmt.Errorf("cannot parse config from environment: %w", errEnv)
	}

	overrideWithCliParams(config)

	return config, nil
}

func overrideWithCliParams(config *Config) {
	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "Database connection string")
	flag.Func("b", "Base URL", func(flagValue string) error {
		if flagValue == "" {
			return nil
		}

		baseURL, errParse := url.Parse(flagValue)
		if errParse != nil {
			return fmt.Errorf("cannot parse [%s] as URL: %w", flagValue, errParse)
		}
		config.BaseURL = *baseURL

		return nil
	})
	flag.Func("f", "DEPRECATED file storage path - not used in favour of DB", func(flagValue string) error {
		if flagValue != "" {
			log.Printf("Detected usage of DEPRECATED flag -f: [%s]", flagValue)
		}

		return nil
	})
	flag.Parse()
}
