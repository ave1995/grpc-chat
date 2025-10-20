package connector

import (
	"context"
	"log"
	"sync"

	"github.com/ave1995/grpc-chat/domain/model"
)

type MessageHub struct {
	subscribers  map[SubscriberID]*MessageSubscriber
	broadcastQue chan *model.Message
	mu           sync.Mutex
}

func NewMessageHub(ctx context.Context, capacity int) *MessageHub {
	h := &MessageHub{
		subscribers:  make(map[SubscriberID]*MessageSubscriber),
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
					close(subscriber.messages)
				}
				h.subscribers = make(map[SubscriberID]*MessageSubscriber) // Proč tohle?
			}()
			return

		case msg := <-h.broadcastQue:
			func() {
				h.mu.Lock()
				defer h.mu.Unlock()
				for _, subscriber := range h.subscribers {
					select {
					case subscriber.messages <- msg:
					default:
						log.Printf("Hub: dropped message for subscriber %v, channel full", subscriber)
					}
				}
			}()
		}
	}
}

func (h *MessageHub) Subscribe(subscriber *MessageSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[subscriber.id] = subscriber
}

func (h *MessageHub) Unsubscribe(subscriber *MessageSubscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscribers, subscriber.id)
}

func (h *MessageHub) Broadcast(msg *model.Message) {
	select {
	case h.broadcastQue <- msg:
	default:
		log.Printf("Hub: broadcast queue full, dropping message: %+v", msg)
	}
}
