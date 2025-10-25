package config

import "time"

type MessageProcessorConfig struct {
	Topic           string
	OutboxInterval  time.Duration
	OutboxBatchSize int
}
