package connector

import "context"

// Consumer defines the interface for reading messages from Kafka topics.
type Consumer interface {
	ReadMessage(ctx context.Context) (topic string, key string, value string, err error)
	Close() error
}
