package server

import (
	"context"
	"log/slog"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/service"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/google/uuid"
)

type ChatServer struct {
	pb.ChatServiceServer
	logger         *slog.Logger
	messageService service.MessageService
}

func NewChatServer(logger *slog.Logger, messageService service.MessageService) *ChatServer {
	return &ChatServer{
		logger:         logger,
		messageService: messageService}
}

// Unary: SendMessage
func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	created, err := s.messageService.SendMessage(ctx, msg.Message.Text)
	if err != nil {
		return nil, err
	}
	s.logger.Info("message sent", "message", created.Text)

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

// Server streaming
func (s *ChatServer) Reader(req *pb.ReaderRequest, srv pb.ChatService_ReaderServer) error {
	s.logger.Info("server stream opened")

	subscriber, cleanup := s.messageService.NewSubscriberWithCleanup()
	defer cleanup()

	for {
		select {
		case <-srv.Context().Done():
			s.logger.Info("server stream closed by disconnection of client: %v", "subscriber", subscriber)
			return nil

		case msg, open := <-subscriber.Messages():
			if !open {
				return nil
			}
			err := srv.Send(msg.ToProto())
			if err != nil {
				s.logger.Error("send error", utils.SlogError(err))
				return err
			}
		}
	}
}
