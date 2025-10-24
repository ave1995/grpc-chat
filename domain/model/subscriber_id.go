package model

import "github.com/google/uuid"

type SubscriberID uuid.UUID

func NewSubscriberID() SubscriberID {
	return SubscriberID(uuid.New())
}
