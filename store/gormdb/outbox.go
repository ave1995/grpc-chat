package gormdb

import (
	"encoding/json"
	"time"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/google/uuid"
)

type outboxEvent struct {
	ID          uuid.UUID          `gorm:"type:uuid;primary_key"`
	EventType   model.EventType    `gorm:"not null"`
	AggregateID uuid.UUID          `gorm:"type:uuid;not null"`
	Payload     json.RawMessage    `gorm:"type:jsonb; not null"`
	Status      model.OutboxStatus `gorm:"not null;default:1"`
	CreatedAt   time.Time          `gorm:"not null;autoCreateTime"`
	ProcessedAt *time.Time         `gorm:"type:timestamp;default:null"`
}

func (e *outboxEvent) toDomain() *model.OutboxEvent {
	return &model.OutboxEvent{
		ID:          model.OutboxEventID(e.ID),
		EventType:   e.EventType,
		Payload:     e.Payload,
		AggregateID: e.AggregateID,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		ProcessedAt: e.ProcessedAt,
	}
}
