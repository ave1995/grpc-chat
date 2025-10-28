package connector

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
)

// Consumer defines the interface for reading messages from Kafka topics.
type Consumer interface {
	Read(ctx context.Context) (<-chan model.Message, error)
}
