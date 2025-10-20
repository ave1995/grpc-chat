package gormdb

import (
	"context"
	"time"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageStore struct {
	gorm *gorm.DB
}

func NewMessageStore(gorm *gorm.DB) store.MessageStore {
	return &messageStore{gorm: gorm}
}

// CreateMessage implements store.MessageStore.
func (m *messageStore) CreateMessage(ctx context.Context, text string) (*model.Message, error) {
	message := &message{
		ID:        uuid.New(),
		Text:      text,
		Timestamp: time.Now(),
	}

	if err := m.gorm.WithContext(ctx).Create(message).Error; err != nil {
		return nil, err
	}

	return message.ToDomain(), nil
}

// GetMessage implements store.MessageStore.
func (m *messageStore) GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error) {
	var message *message
	if err := m.gorm.First(&message, "id = ?", uuid.UUID(id)).Error; err != nil {
		return nil, err
	}

	return message.ToDomain(), nil
}
