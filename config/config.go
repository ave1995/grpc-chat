package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"DB_HOST"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	DBName   string `envconfig:"DB_NAME"`
	Port     string `envconfig:"DB_PORT"`
	Brokers  string `envconfig:"KAFKA_BROKERS"`
}

const EnvPrefix = ""

func NewConfig() (Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

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

func (c Config) KafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: strings.Split(c.Brokers, ","),
	}
}
