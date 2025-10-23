package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host                             string `envconfig:"DB_HOST"`
	User                             string `envconfig:"DB_USER"`
	Password                         string `envconfig:"DB_PASSWORD"`
	DBName                           string `envconfig:"DB_NAME"`
	Port                             string `envconfig:"DB_PORT"`
	Brokers                          string `envconfig:"KAFKA_BROKERS"`
	MessageServiceTopic              string `envconfig:"MESSAGE_SERVICE_TOPIC"`
	MessageServiceSubscriberCapacity int    `envconfig:"MESSAGE_SERVICE_SUBSCRIBER_CAPACITY"`
	HubCapacity                      int    `envconfig:"MESSAGE_HUB_CAPICITY"`
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

func (c Config) KafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: strings.Split(c.Brokers, ","),
	}
}

func (c Config) MessageServiceConfig() MessageServiceConfig {
	return MessageServiceConfig{
		Topic:              c.MessageServiceTopic,
		SubscriberCapacity: c.MessageServiceSubscriberCapacity,
	}
}

func (c Config) HubConfig() HubConfig {
	return HubConfig{
		Capicity: c.HubCapacity,
	}
}
