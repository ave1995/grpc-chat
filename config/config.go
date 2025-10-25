package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host                             string        `envconfig:"DB_HOST"`
	User                             string        `envconfig:"DB_USER"`
	Password                         string        `envconfig:"DB_PASSWORD"`
	DBName                           string        `envconfig:"DB_NAME"`
	Port                             string        `envconfig:"DB_PORT"`
	Brokers                          string        `envconfig:"KAFKA_BROKERS"`
	MessageTopic                     string        `envconfig:"MESSAGE_TOPIC"`
	MessageServiceSubscriberCapacity int           `envconfig:"MESSAGE_SERVICE_SUBSCRIBER_CAPACITY"`
	HubCapacity                      int           `envconfig:"MESSAGE_HUB_CAPACITY"`
	OutboxInterval                   time.Duration `envconfig:"OUTBOX_INTERVAL" default:"1s"`
	OutboxBatchSize                  int           `envconfig:"OUTBOX_BATCH_SIZE" default:"100"`
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
		SubscriberCapacity: c.MessageServiceSubscriberCapacity,
	}
}

func (c Config) HubConfig() HubConfig {
	return HubConfig{
		Capacity: c.HubCapacity,
	}
}

func (c Config) MessageProcessorConfig() MessageProcessorConfig {
	return MessageProcessorConfig{
		Topic:           c.MessageTopic,
		OutboxInterval:  c.OutboxInterval,
		OutboxBatchSize: c.OutboxBatchSize,
	}
}
