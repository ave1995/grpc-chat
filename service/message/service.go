package message

import (
	"context"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
)

var _ service.MessageService = (*messageService)(nil)

type messageService struct {
	config     config.MessageServiceConfig
	store      store.MessageStore
	messageHub *MessageHub
	producer   connector.Producer
}

func NewMessageService(config config.MessageServiceConfig, store store.MessageStore, producer connector.Producer, messageHub *MessageHub) service.MessageService {
	return &messageService{
		config:     config,
		store:      store,
		messageHub: messageHub,
		producer:   producer,
	}
}

func (m *messageService) GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error) {
	return m.store.GetMessage(ctx, id)
}

// TODO: outbox transactional pattern
func (m *messageService) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	msg, err := m.store.CreateMessage(ctx, text)
	if err != nil {
		return nil, err
	}

	err = m.producer.SendMessage(ctx, m.config.Topic, msg.ID.String(), msg.Text)
	if err != nil {
		return nil, err
	}

	m.messageHub.Broadcast(msg)

	return msg, nil
}

func (m *messageService) NewSubscriberWithCleanup() (*model.MessageSubscriber, func()) {
	subscriber := model.NewSubscriber(model.SubscriberID(uuid.New()), m.config.SubscriberCapacity)

	m.messageHub.Subscribe(subscriber)
	cleanup := func() { m.messageHub.Unsubscribe(subscriber) }

	return subscriber, cleanup
}
