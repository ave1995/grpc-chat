package message

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
)

type messageService struct {
	store store.MessageStore
}

func NewMessageService(store store.MessageStore) service.MessageService {
	return &messageService{store: store}
}

// GetMessage implements service.MessageService.
func (m *messageService) GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	return m.store.GetMessage(ctx, id)
}

// SendMessage implements service.MessageService.
func (m *messageService) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	return m.store.CreateMessage(ctx, text)
}
