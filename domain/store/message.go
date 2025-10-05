package store

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/google/uuid"
)

type MessageStore interface {
	GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error)
	CreateMessage(ctx context.Context, text string) (*model.Message, error)
}
