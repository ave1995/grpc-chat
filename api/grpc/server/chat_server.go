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

// Unary: SendMessage
func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	created, err := s.messageService.SendMessage(ctx, msg.Message.Text)
	if err != nil {
		return nil, err
	}

	s.messageHub.Broadcast(created)

	log.Printf("message sent: %s", created.Text)
	// TODO: better response
	return &pb.SendMessageResponse{Message: "Message stored successfully", Id: created.ID.String()}, nil
}

// Unary: GetMessage
func (s *ChatServer) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	u, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	id := model.MessageID(u)

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

// Bidirectional stream: Chat
func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("chat stream opened")

	ctx := stream.Context()

	subscriber := connector.NewSubscriber(connector.SubscriberID(uuid.New()), 10)

	s.messageHub.Subscribe(subscriber)
	defer s.messageHub.Unsubscribe(subscriber)

	log.Printf("Client subscribed: %v", subscriber)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Client disconnected: %v", subscriber)
			return nil

		case msg, ok := <-subscriber.Messages():
			if !ok {
				return nil // hub closed
			}
			// Convert model.Message to protobuf message and send
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
