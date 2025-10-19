package connector

import (
	"context"
	"log"
	"sync"

	"github.com/ave1995/grpc-chat/domain/model"
)

type Subscriber chan model.Message

type MessageHub struct {
	subscribers  map[Subscriber]any
	broadcastQue chan model.Message
	mu           sync.Mutex
	ctx          context.Context
}

func NewMessageHub(ctx context.Context, capacity int) *MessageHub {
	h := &MessageHub{
		subscribers:  make(map[Subscriber]any),
		broadcastQue: make(chan model.Message, capacity),
		ctx:          ctx,
	}
	go h.run()
	return h
}

func (h *MessageHub) run() {
	for {
		select {
		case <-h.ctx.Done():
			h.mu.Lock()
			for subscriber := range h.subscribers {
				close(subscriber)
			}
			h.mu.Unlock()
			return

		case msg := <-h.broadcastQue:
			h.mu.Lock()
			subs := make([]Subscriber, 0, len(h.subscribers))
			for subscriber := range h.subscribers {
				subs = append(subs, subscriber)
			}
			h.mu.Unlock()

			for _, subscriber := range subs {
				select {
				case subscriber <- msg:
				default:
					log.Printf("Hub: dropped message for subscriber %p, channel full", subscriber)
				}
			}
		}
	}
}

func (h *MessageHub) Subscribe(subscriber Subscriber) {
	h.mu.Lock()
	h.subscribers[subscriber] = nil
	h.mu.Unlock()
}

func (h *MessageHub) Unsubscribe(subscriber Subscriber) {
	h.mu.Lock()
	delete(h.subscribers, subscriber)
	h.mu.Unlock()
}

func (h *MessageHub) Broadcast(msg model.Message) {
	h.broadcastQue <- msg
}
