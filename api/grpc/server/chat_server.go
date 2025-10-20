package server

import (
	"context"
	"log"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/google/uuid"
)

type ChatServer struct {
	pb.ChatServiceServer
	messageService service.MessageService
	messageHub     *connector.MessageHub
}

func NewChatServer(messageService service.MessageService, messageHub *connector.MessageHub) *ChatServer {
	return &ChatServer{messageService: messageService, messageHub: messageHub}
}

// špatně komentář
// Unary: SendMessage
func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	created, err := s.messageService.SendMessage(ctx, msg.Message.Text)
	if err != nil {
		return nil, err
	}

	// Tohle tady nemá co dělat, to má být v service.
	s.messageHub.Broadcast(created)

	log.Printf("message sent: %s", created.Text) // co tenhle log?
	// TODO: better response co znamená toto todo?
	return &pb.SendMessageResponse{Message: "Message stored successfully", Id: created.ID.String()}, nil
}

// Unary: GetMessage
func (s *ChatServer) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	// Co kdyby si tohle celé přesunul do modelu pod to MessageID jako nějakou funkci, ta by se mohla hodit na více místech.
	u, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	id := model.MessageID(u) // tohle můžeš udělat v argumentu te metody

	// found zní jak boolean jestli to bylo nalezené nebo ne....
	found, err := s.messageService.GetMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetMessageResponse{
		Message: &pb.Message{
			Text: found.Text,
		},
	}, nil
}

// Bidirectional stream: Chat špatně komentář
func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("chat stream opened") // špatně log

	ctx := stream.Context() // na co je tahle proměnná, když jí používáš na jednom místě?

	// Tady tu logiku bych částečně taky schoval za service.
	// Tenhle subscriber bych tu nechal. Ale Registrace/Deregistrace bude schovaná za service.
	subscriber := connector.NewSubscriber(connector.SubscriberID(uuid.New()), 10)

	s.messageHub.Subscribe(subscriber)
	defer s.messageHub.Unsubscribe(subscriber)

	log.Printf("Client subscribed: %v", subscriber) // oprav všude ten log na slog a jen jednu instanci

	for {
		select {
		case <-ctx.Done():
			log.Printf("Client disconnected: %v", subscriber)
			return nil

		case msg, ok := <-subscriber.Messages():
			if !ok {
				// tohle nikdy nemůže nastat ale pro uvolnění paměti by možná mělo iděálně udělat Close funkci na Subscriber a tam channel zavřít
				// Ale musíš dát pozor aby to nešlo zavřít víckrát a aby se ti nestalo, že zavřeš a potom zapíšeš do channelu = panika
				// Ale takhle je to taky OK. Kdyby si to předělával tak to nepřekombinuj zase :D.
				// Dodatečně jsem našel místo, kde to zavíráš, když vyprší kontext. Měj stejný přístup na odhlašování těch subscriberů,
				// jak v případě vypnutí aplikace, tak když se odhlásí klient.
				return nil // hub closed
			}
			// Convert model.Message to protobuf message and send
			// Na vytvoření pb z modelu využíváme většinou separe package (můžu ti potom ukázat jak a nebo klidně na modelu mít funkci na to)
			err := stream.Send(&pb.Message{
				Text: msg.Text,
			})
			if err != nil {
				log.Printf("Send error: %v", err)
				return err
			}
		}
	}
}
