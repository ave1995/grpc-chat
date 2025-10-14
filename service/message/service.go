package message

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
)

type messageService struct {
	store    store.MessageStore
	producer store.Producer
	topic    string
}

func NewMessageService(store store.MessageStore, producer store.Producer, topic string) service.MessageService {
	return &messageService{
		store:    store,
		producer: producer,
		topic:    topic,
	}
}

// GetMessage implements service.MessageService.
func (m *messageService) GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	return m.store.GetMessage(ctx, id)
}

// SendMessage implements service.MessageService.
func (m *messageService) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	// 1. Create message in DB
	msg, err := m.store.CreateMessage(ctx, text)
	if err != nil {
		return nil, err
	}

	// 2. Publish message event to Kafka (non-blocking)
	go func() {
		payload, _ := json.Marshal(msg)
		err := m.producer.SendMessage(ctx, m.topic, msg.ID.String(), string(payload))
		if err != nil {
			log.Printf("failed to send Kafka message: %v", err)
		}
	}()

	return msg, nil
}
