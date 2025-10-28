package message

import (
	"context"
	"log/slog"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/google/uuid"
)

var _ service.MessageService = (*Service)(nil)

type Service struct {
	logger     *slog.Logger
	config     config.MessageServiceConfig
	store      store.MessageStore
	messageHub *Hub
	consumer   connector.Consumer
}

func NewService(logger *slog.Logger, config config.MessageServiceConfig, store store.MessageStore, messageHub *Hub, consumer connector.Consumer) *Service {
	return &Service{
		logger:     logger,
		config:     config,
		store:      store,
		messageHub: messageHub,
		consumer:   consumer,
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

func (m *Service) Broadcast(ctx context.Context) {
	go func() {
		msgCh, err := m.consumer.Read(ctx)
		if err != nil {
			m.logger.Error("consumer.read", utils.SlogError(err))
		}
		for msg := range msgCh {
			err := m.messageHub.Broadcast(&msg)
			if err != nil {
				m.logger.Error("hub.broadcast", utils.SlogError(err))
			}
		}
	}()
}
