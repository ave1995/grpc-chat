package connector

import (
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/google/uuid"
)

// Otázka, že když tyhle structy používáš mimo package tak jestli by neměly být v modelu.
type SubscriberID uuid.UUID

type MessageSubscriber struct {
	id       SubscriberID
	messages chan *model.Message
}

func NewSubscriber(id SubscriberID, capacity int) *MessageSubscriber {
	return &MessageSubscriber{
		id:       id,
		messages: make(chan *model.Message, capacity),
	}
}

func (s *MessageSubscriber) ID() SubscriberID {
	return s.id
}

func (s *MessageSubscriber) Messages() <-chan *model.Message {
	return s.messages
}
