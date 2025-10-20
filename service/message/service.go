package message

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/domain/store"
)

type messageService struct {
	store    store.MessageStore
	producer connector.Producer
	topic    string
}

// Nauč se používat tento zápis: tím explicitne vynutíš použití interface.
//var _ service.MessageService = (*messageService)(nil)

func NewMessageService(store store.MessageStore, producer connector.Producer, topic string) service.MessageService {
	return &messageService{
		store:    store,
		producer: producer,
		topic:    topic,
	}
}

func (m *messageService) GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error) {
	return m.store.GetMessage(ctx, id) // vrapuj errory
}

// TODO: outbox transactional pattern
func (m *messageService) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	msg, err := m.store.CreateMessage(ctx, text) // radši save nebo store message.
	if err != nil {
		return nil, err
	}

	// Ty máš tady správně, že service odešle do kafky a potom v GRPC to pošleš do toho HUBU :D
	// S tím, že consumer jsem nenašel použitý, takže to je takový zvláštní hybrid.
	// Proč máš topic v v service jako atribut, když ho tu nepoužijes?
	err = m.producer.SendMessage(ctx, "messages", msg.ID.String(), msg.Text)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
