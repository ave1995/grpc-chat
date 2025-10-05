package gormdb

import (
	"time"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/google/uuid"
)

type message struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Text      string
	Timestamp time.Time
}

func (m *message) ToDomain() *model.Message {
	return &model.Message{
		ID:        m.ID,
		Text:      m.Text,
		Timestamp: m.Timestamp,
	}
}
