package store

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
)

type MessageStore interface {
	GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error)
	CreateMessage(ctx context.Context, text string) (*model.Message, error)
}
