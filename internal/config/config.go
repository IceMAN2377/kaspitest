package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	PostgresMigrate  bool   `env:"POSTGRES_MIGRATE" envDefault:"true"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDb       string `env:"POSTGRES_DB,required"`
	PostgresSslMode  string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`

	HttpPort int `env:"HTTP_PORT" envDefault:"8080"`
}

// NewConfig creates a new Config
func NewConfig() *Config {
	cfg := &Config{}
	if err := cfg.readFromEnvironment(); err != nil {
		panic(err)
	}
	return cfg
}

// readFromEnvironment reads the settings from environment variables.
func (c *Config) readFromEnvironment() error {
	return env.Parse(c)
}
