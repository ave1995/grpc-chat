package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID          uuid.UUID
	EventType   EventType
	Payload     json.RawMessage
	Status      OutboxStatus
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
