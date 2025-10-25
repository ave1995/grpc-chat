package gormdb

import (
	"context"
	"encoding/json"
	"errors"

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

func (m *MessageStore) Create(ctx context.Context, text string) (*model.Message, error) {
	var msg *message

	err := m.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		msg = &message{
			ID:   uuid.New(),
			Text: text,
		}

		if err := m.gorm.WithContext(ctx).Create(msg).Error; err != nil {
			return err
		}

		payloadBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		outbox := &outboxEvent{
			ID:        uuid.New(),
			EventType: model.SendMessage,
			Payload:   payloadBytes,
			Status:    model.Pending,
		}

		if err := m.gorm.WithContext(ctx).Create(outbox).Error; err != nil {
			return err
		}

		return nil
	})

	return msg.ToDomain(), err
}

func (m *MessageStore) Fetch(ctx context.Context, id model.MessageID) (*model.Message, error) {
	var msg *message
	if err := m.gorm.WithContext(ctx).First(&msg, "id = ?", uuid.UUID(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return msg.ToDomain(), nil
}
