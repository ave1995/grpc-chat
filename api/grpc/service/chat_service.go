package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	mu       sync.Mutex
	messages map[string]*pb.Message
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		messages: make(map[string]*pb.Message),
	}
}

// Unary: SendMessage
func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg.Timestamp = time.Now().Unix()
	s.messages[msg.Id] = msg

	log.Printf("Received message: %s", msg.Text)
	return &pb.Ack{Message: "Message received"}, nil
}

// Unary: GetMessage
func (s *ChatServer) GetMessage(ctx context.Context, req *pb.MessageRequest) (*pb.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[req.Id]
	if !ok {
		return nil, fmt.Errorf("message with ID %s not found", req.Id)
	}

	return msg, nil
}

// Bidirectional stream: Chat
func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		s.mu.Lock()
		msg.Timestamp = time.Now().Unix()
		s.messages[msg.Id] = msg
		s.mu.Unlock()

		log.Printf("Streaming message: %s", msg.Text)

		if err := stream.Send(msg); err != nil {
			return err
		}
	}
}
