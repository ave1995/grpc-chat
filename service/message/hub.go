package message

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/segmentio/kafka-go"
)

var _ connector.Broadcaster = (*Hub)(nil)

type Hub struct {
	logger       *slog.Logger
	subscribers  map[model.SubscriberID]*model.MessageSubscriber
	broadcastQue chan *model.Message
	mu           sync.Mutex
}

func NewHub(ctx context.Context, logger *slog.Logger, capacity int) *Hub {
	h := &Hub{
		logger:       logger,
		subscribers:  make(map[model.SubscriberID]*model.MessageSubscriber),
		broadcastQue: make(chan *model.Message, capacity),
	}
	go h.run(ctx)
	return h
}

func (h *Hub) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			func() {
				h.mu.Lock()
				defer h.mu.Unlock()
				for _, subscriber := range h.subscribers {
					subscriber.Close()
				}
			}()
			return

		case msg := <-h.broadcastQue:
			func() {
				h.mu.Lock()
				defer h.mu.Unlock()
				for _, subscriber := range h.subscribers {
					subscriber.Push(msg)
				}
			}()
		}
	}
}

func (h *Hub) Subscribe(subscriber *model.MessageSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[subscriber.ID()] = subscriber
}

func (h *Hub) Unsubscribe(subscriber *model.MessageSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscribers, subscriber.ID())
}

func (h *Hub) Broadcast(event kafka.Message) error {
	msg := &model.Message{
		ID:   model.MessageID(event.Key),
		Text: string(event.Value),
	}

	select {
	case h.broadcastQue <- msg:
	default:
		h.logger.Warn("Hub: dropped message for broadcast, channel full", "message", event)
		return errors.New("hub: broadcast queue full")
	}

	return nil
}
