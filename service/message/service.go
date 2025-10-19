package message

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
)

type messageService struct {
	store    store.MessageStore
	producer connector.Producer
	topic    string
}

func NewMessageService(store store.MessageStore, producer connector.Producer, topic string) service.MessageService {
	return &messageService{
		store:    store,
		producer: producer,
		topic:    topic,
	}
}

func (m *messageService) GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	return m.store.GetMessage(ctx, id)
}

// TODO: outbox transactional pattern
func (m *messageService) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	msg, err := m.store.CreateMessage(ctx, text)
	if err != nil {
		return nil, err
	}

	err = m.producer.SendMessage(ctx, "messages", msg.ID.String(), msg.Text)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
