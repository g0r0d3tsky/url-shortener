package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Postgres struct {
		Host     string `env:"POSTGRES_HOST,notEmpty"`
		Port     string `env:"POSTGRES_PORT,notEmpty"`
		User     string `env:"POSTGRES_USER,notEmpty"`
		Password string `env:"POSTGRES_PASSWORD,notEmpty"`
		Database string `env:"POSTGRES_DB,notEmpty"`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST,notEmpty"`
		Port     string `env:"REDIS_PORT,notEmpty"`
		Database string `env:"REDIS_DB,notEmpty"`
		Key      string `env:"REDIS_KEY"`
	}
	MonthAmount int `env:"MONTH_AMOUNT,notEmpty"`
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Database,
	)
}

func (c *Config) RedisDSN() string {
	return fmt.Sprintf("redis://%s:%s/%s",
		c.Redis.Host, c.Redis.Port, c.Redis.Database,
	)
}

func Read() (*Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}
