package service

import (
	"context"

	"github.com/ave1995/grpc-chat/connector"
	"github.com/ave1995/grpc-chat/domain/model"
)

type MessageService interface {
	SendMessage(ctx context.Context, text string) (*model.Message, error)
	GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error)
	NewSubscriberWithCleanup() (*connector.MessageSubscriber, func())
}
