package model

type MessageSubscriber struct {
	id       SubscriberID
	messages chan *Message
}

func NewSubscriber(id SubscriberID, capacity int) *MessageSubscriber {
	return &MessageSubscriber{
		id:       id,
		messages: make(chan *Message, capacity),
	}
}

func (s *MessageSubscriber) ID() SubscriberID {
	return s.id
}

func (s *MessageSubscriber) Messages() <-chan *Message {
	return s.messages
}

func (s *MessageSubscriber) Close() {
	close(s.messages)
}

func (s *MessageSubscriber) Push(msg *Message) bool {
	select {
	case s.messages <- msg:
		return true
	default:
		return false
	}
}
