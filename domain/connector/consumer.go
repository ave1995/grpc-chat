package connector

import "context"

// Consumer defines the interface for reading messages from Kafka topics.
type Consumer interface {
	Start(ctx context.Context)
}
