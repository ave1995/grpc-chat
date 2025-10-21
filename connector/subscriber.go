package connector

import (
	"github.com/ave1995/grpc-chat/domain/model"
)

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
