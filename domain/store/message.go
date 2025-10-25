package store

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
)

type MessageStore interface {
	Fetch(ctx context.Context, id model.MessageID) (*model.Message, error)
	Create(ctx context.Context, text string) (*model.Message, error)
}
