package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Addr string   `env:"ADDRESS" envDefault:"localhost:8080"`
	Db   DbConfig // Embedded struct for DB configuration
	AUTH0_DOMAIN string `env:"AUTH0_DOMAIN" envDefault:"https://dev-n5mocwlrk8i63cjm.us.auth0.com/"`
	AUTH0_AUDIENCE string `env:"AUTH0_AUDIENCE" envDefault:"https://project-management-api"`
	ENVIRONMENT string	`env:"ENVIRONMENT" envDefault:"local"`
}

type DbConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"password"`
	DBName   string `env:"DB_NAME" envDefault:"project_management_system"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

func LoadEnvConfigs() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %+v", err)
	}
	return &cfg
}
