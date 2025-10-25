package connector

import "context"

// Producer defines the interface for sending messages to Kafka topics.
type Producer interface {
	Send(ctx context.Context, topic string, key string, value []byte) error
	Close() error
}
