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

var _ store.MessageStore = (*messageStore)(nil)

type messageStore struct {
	gorm *gorm.DB
}

func NewMessageStore(gorm *gorm.DB) store.MessageStore {
	return &messageStore{gorm: gorm}
}

func (m *messageStore) CreateMessage(ctx context.Context, text string) (*model.Message, error) {
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

func (m *messageStore) GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error) {
	var msg *message
	if err := m.gorm.WithContext(ctx).First(&msg, "id = ?", uuid.UUID(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return msg.ToDomain(), nil
}
