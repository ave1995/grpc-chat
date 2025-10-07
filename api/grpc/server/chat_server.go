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
func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	created, err := s.messageService.SendMessage(ctx, msg.Message.Text)
	if err != nil {
		return nil, err
	}

	log.Printf("message sent: %s", created.Text)
	// TODO: better response
	return &pb.SendMessageResponse{Message: "Message stored successfully", Id: created.ID.String()}, nil
}

// Unary: GetMessage
func (s *ChatServer) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

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

	for {
		// Receive message from client
		msg, err := stream.Recv()
		if err != nil {
			// client closed stream
			log.Printf("chat stream closed: %v", err)
			return err
		}

		log.Printf("received message: %s", msg.Text)

		// Store message — your service sets timestamp etc.
		stored, err := s.messageService.SendMessage(stream.Context(), msg.Text)
		if err != nil {
			log.Printf("failed to store message: %v", err)
			continue
		}

		// Echo back the stored version (now with ID + timestamp)
		if err := stream.Send(&pb.Message{
			Text: stored.Text,
		}); err != nil {
			log.Printf("failed to send back: %v", err)
			return err
		}
	}
}
