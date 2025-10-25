package gormdb

import (
	"encoding/json"
	"time"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/google/uuid"
)

type outboxEvent struct {
	ID        uuid.UUID          `gorm:"type:uuid;primary_key"`
	EventType model.EventType    `gorm:"not null"`
	Payload   json.RawMessage    `gorm:"type:jsonb; not null"`
	Status    model.OutboxStatus `gorm:"not null;default:1"`
	CreatedAt time.Time          `gorm:"not null;autoCreateTime"`
	ProceedAt time.Time          `gorm:"autoUpdateTime"`
}
