package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTPServer HTTPServer
		Postgres   Postgres
	}

	HTTPServer struct {
		GinEnviroment          string `env:"GIN_ENV" envDefault:"debug"`
		Port                   int    `env:"HTTP_PORT,notEmpty"`
		ShutdownTimeoutSeconds int    `env:"HTTP_SHUTDOWN_TIMEOUT_SECONDS" envDefault:"10"`
	}

	Postgres struct {
		Host         string        `env:"POSTGRES_HOST,notEmpty"`
		Port         string        `env:"POSTGRES_PORT,notEmpty"`
		User         string        `env:"POSTGRES_USER,notEmpty"`
		Password     string        `env:"POSTGRES_PASSWORD,notEmpty"`
		DBName       string        `env:"POSTGRES_DB,notEmpty"`
		MaxPoolSize  int           `env:"POSTGRES_MAX_POOL_SIZE" envDefault:"10"`
		ConnAttempts int           `env:"POSTGRES_CONN_ATTEMPTS" envDefault:"5"`
		ConnTimeout  time.Duration `env:"POSTGRES_CONN_TIMEOUT" envDefault:"5s"`
	}
)

func (c Postgres) GetDsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)
}

func New(path string) (Config, error) {
	var config Config

	if err := godotenv.Load(path); err != nil {
		return config, fmt.Errorf("failed to load config: %w", err)
	}

	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}
