package connector

import "context"

// Producer defines the interface for sending messages to Kafka topics.
type Producer interface {
	SendMessage(ctx context.Context, topic string, key string, value string) error
	Close() error
}
