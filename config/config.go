package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

const EnvPrefix = ""

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process(EnvPrefix, &cfg); err != nil {
		return Config{}, fmt.Errorf("procces new config: %w", err)
	}
	return cfg, nil
}

func (c Config) DBConfig() DBConfig {
	return DBConfig{
		Host:     c.Host,
		User:     c.User,
		Password: c.Password,
		DBName:   c.DBName,
		Port:     c.Port,
	}
}
