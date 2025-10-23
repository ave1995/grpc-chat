package message

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ave1995/grpc-chat/domain/model"
)

type MessageHub struct {
	logger       *slog.Logger
	subscribers  map[model.SubscriberID]*model.MessageSubscriber
	broadcastQue chan *model.Message
	mu           sync.Mutex
}

func NewMessageHub(ctx context.Context, logger *slog.Logger, capacity int) *MessageHub {
	h := &MessageHub{
		logger:       logger,
		subscribers:  make(map[model.SubscriberID]*model.MessageSubscriber),
		broadcastQue: make(chan *model.Message, capacity),
	}
	go h.run(ctx)
	return h
}

func (h *MessageHub) run(ctx context.Context) {
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

func (h *MessageHub) Subscribe(subscriber *model.MessageSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[subscriber.ID()] = subscriber
}

func (h *MessageHub) Unsubscribe(subscriber *model.MessageSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscribers, subscriber.ID())
}

func (h *MessageHub) Broadcast(msg *model.Message) {
	select {
	case h.broadcastQue <- msg:
	default:
		h.logger.Warn("Hub: dropped message for broadcast, channel full", "message", msg)
	}
}
