package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID          OutboxEventID
	EventType   EventType
	AggregateID uuid.UUID
	Payload     json.RawMessage
	Status      OutboxStatus
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
