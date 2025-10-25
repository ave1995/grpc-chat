package message

import (
	"context"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
)

var _ service.MessageService = (*Service)(nil)

type Service struct {
	config     config.MessageServiceConfig
	store      store.MessageStore
	messageHub *Hub
}

func NewService(config config.MessageServiceConfig, store store.MessageStore, messageHub *Hub) *Service {
	return &Service{
		config:     config,
		store:      store,
		messageHub: messageHub,
	}
}

func (m *Service) Fetch(ctx context.Context, id model.MessageID) (*model.Message, error) {
	return m.store.Fetch(ctx, id)
}

func (m *Service) Send(ctx context.Context, text string) (*model.Message, error) {
	msg, err := m.store.Create(ctx, text)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Service) NewSubscriberWithCleanup() (*model.MessageSubscriber, func()) {
	subscriber := model.NewSubscriber(model.SubscriberID(uuid.New()), m.config.SubscriberCapacity)

	m.messageHub.Subscribe(subscriber)
	cleanup := func() { m.messageHub.Unsubscribe(subscriber) }

	return subscriber, cleanup
}
