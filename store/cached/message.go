package cached

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/ave1995/grpc-chat/utils"
)

var _ store.MessageStore = (*MessageStore)(nil)

type MessageStore struct {
	logger *slog.Logger
	store  store.MessageStore
	cache  store.Cache
}

func NewMessageStore(logger *slog.Logger, store store.MessageStore, cache store.Cache) *MessageStore {
	return &MessageStore{logger: logger, store: store, cache: cache}
}

func (c MessageStore) Fetch(ctx context.Context, id model.MessageID) (*model.Message, error) {
	val, ok, err := c.cache.Get(ctx, id.String())
	if ok {
		c.logger.Info("→ cache hit")
		return val.(*model.Message), nil
	}
	if err != nil {
		c.logger.Error("→ cache error", utils.SlogError(err))
	}

	c.logger.Info("→ cache miss", "id", id)

	msg, err := c.store.Fetch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cache fetch error: %w", err)
	}

	err = c.cache.Set(ctx, msg.ID.String(), msg, 0)
	if err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return msg, nil
}

func (c MessageStore) Create(ctx context.Context, text string) (*model.Message, error) {
	msg, err := c.store.Create(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("cache create error: %w", err)
	}

	err = c.cache.Set(ctx, msg.ID.String(), msg, 0)
	if err != nil {
		return nil, fmt.Errorf("cache set error: %w", err)
	}

	return msg, nil
}
