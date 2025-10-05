package service

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/google/uuid"
)

type MessageService interface {
	SendMessage(ctx context.Context, text string) (*model.Message, error)
	GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error)
}
