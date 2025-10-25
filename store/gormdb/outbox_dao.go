package gormdb

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/ave1995/grpc-chat/utils"
	"gorm.io/gorm"
)

var _ store.OutboxStore = (*OutboxStore)(nil)

type OutboxStore struct {
	gorm *gorm.DB
}

func NewOutboxStore(gorm *gorm.DB) *OutboxStore {
	return &OutboxStore{gorm: gorm}
}

func (o *OutboxStore) GetPendingEvents(ctx context.Context, eType model.EventType, limit int) ([]*model.OutboxEvent, error) {
	var rows []outboxEvent
	err := o.gorm.WithContext(ctx).
		Where(&outboxEvent{Status: model.Pending, EventType: eType}).
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	events := make([]*model.OutboxEvent, len(rows))
	for i, row := range rows {
		events[i] = row.toDomain()
	}
	return events, nil
}

func (o *OutboxStore) MarkProcessed(ctx context.Context, id model.OutboxEventID) error {
	return o.gorm.WithContext(ctx).
		Model(outboxEvent{}).Where("id = ?", id).
		Updates(outboxEvent{Status: model.Processed, ProcessedAt: utils.NowPtr()}).Error
}

func (o *OutboxStore) MarkFailed(ctx context.Context, id model.OutboxEventID) error {
	return o.gorm.WithContext(ctx).
		Model(outboxEvent{}).Where("id = ?", id).
		Updates(outboxEvent{Status: model.Failed}).Error
}
