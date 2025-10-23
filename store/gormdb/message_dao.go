package gormdb

import (
	"context"
	"errors"
	"time"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ store.MessageStore = (*MessageStore)(nil)

type MessageStore struct {
	gorm *gorm.DB
}

func NewMessageStore(gorm *gorm.DB) *MessageStore {
	return &MessageStore{gorm: gorm}
}

func (m *MessageStore) CreateMessage(ctx context.Context, text string) (*model.Message, error) {
	msg := &message{
		ID:        uuid.New(),
		Text:      text,
		Timestamp: time.Now(),
	}

	if err := m.gorm.WithContext(ctx).Create(msg).Error; err != nil {
		return nil, err
	}

	return msg.ToDomain(), nil
}

func (m *MessageStore) GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error) {
	var msg *message
	if err := m.gorm.WithContext(ctx).First(&msg, "id = ?", uuid.UUID(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return msg.ToDomain(), nil
}
