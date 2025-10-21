package message

import (
	"context"

	"github.com/ave1995/grpc-chat/connector"
	cd "github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
)

type messageService struct {
	store      store.MessageStore
	messageHub *connector.MessageHub
	producer   cd.Producer
	topic      string
}

func NewMessageService(store store.MessageStore, producer cd.Producer, topic string, messageHub *connector.MessageHub) service.MessageService {
	return &messageService{
		store:      store,
		messageHub: messageHub,
		producer:   producer,
		topic:      topic,
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

	err = m.producer.SendMessage(ctx, "messages", msg.ID.String(), msg.Text)
	if err != nil {
		return nil, err
	}

	m.messageHub.Broadcast(msg)

	return msg, nil
}

func (m *messageService) NewSubscriberWithCleanup() (*connector.MessageSubscriber, func()) {
	subscriber := connector.NewSubscriber(connector.SubscriberID(uuid.New()), 10)

	m.messageHub.Subscribe(subscriber)
	cleanup := func() { m.messageHub.Unsubscribe(subscriber) }

	return subscriber, cleanup
}
