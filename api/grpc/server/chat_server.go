package server

import (
	"context"
	"log"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/google/uuid"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	messageService service.MessageService
}

func NewChatServer(messageService service.MessageService) *ChatServer {
	return &ChatServer{messageService: messageService}
}

// Unary: SendMessage
func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Ack, error) {
	created, err := s.messageService.SendMessage(ctx, msg.Text)
	if err != nil {
		return nil, err
	}

	log.Printf("Message sent: %s", created.Text)
	return &pb.Ack{Message: "Message stored successfully"}, nil
}

// Unary: GetMessage
func (s *ChatServer) GetMessage(ctx context.Context, req *pb.MessageRequest) (*pb.Message, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	found, err := s.messageService.GetMessage(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.Message{
		Id:        found.ID.String(),
		Text:      found.Text,
		Timestamp: found.Timestamp.Unix(),
	}, nil
}

// Bidirectional stream: Chat
func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("Chat stream opened")

	for {
		// Receive message from client
		msg, err := stream.Recv()
		if err != nil {
			// client closed stream
			log.Printf("Chat stream closed: %v", err)
			return err
		}

		log.Printf("Received message: %s", msg.Text)

		// Store message — your service sets timestamp etc.
		stored, err := s.messageService.SendMessage(stream.Context(), msg.Text)
		if err != nil {
			log.Printf("failed to store message: %v", err)
			continue
		}

		// Echo back the stored version (now with ID + timestamp)
		if err := stream.Send(&pb.Message{
			Id:        stored.ID.String(),
			Text:      stored.Text,
			Timestamp: stored.Timestamp.Unix(),
		}); err != nil {
			log.Printf("failed to send back: %v", err)
			return err
		}
	}
}
