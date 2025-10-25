package server_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/api/grpc/server"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

type serverStreamingMock struct {
	send       func(res *pb.Message) error
	setHeader  func(md metadata.MD) error
	sendHeader func(md metadata.MD) error
	setTrailer func(md metadata.MD)
	context    func() context.Context
	sendMsg    func(m any) error
	recvMsg    func(m any) error
}

func (s *serverStreamingMock) Send(res *pb.Message) error {
	return s.send(res)
}

func (s *serverStreamingMock) SetHeader(md metadata.MD) error {
	return s.setHeader(md)
}

func (s *serverStreamingMock) SendHeader(md metadata.MD) error {
	return s.sendHeader(md)
}

func (s *serverStreamingMock) SetTrailer(md metadata.MD) {
	return
}

func (s *serverStreamingMock) Context() context.Context {
	return s.context()
}

func (s *serverStreamingMock) SendMsg(m any) error {
	return s.sendMsg(m)
}

func (s *serverStreamingMock) RecvMsg(m any) error {
	return s.recvMsg(m)
}

type messageServiceMock struct {
	sendMessage              func(ctx context.Context, text string) (*model.Message, error)
	getMessage               func(ctx context.Context, id model.MessageID) (*model.Message, error)
	newSubscriberWithCleanup func() (*model.MessageSubscriber, func())
}

func (m *messageServiceMock) SendMessage(ctx context.Context, text string) (*model.Message, error) {
	return m.sendMessage(ctx, text)
}

func (m *messageServiceMock) GetMessage(ctx context.Context, id model.MessageID) (*model.Message, error) {
	return m.getMessage(ctx, id)
}

func (m *messageServiceMock) NewSubscriberWithCleanup() (*model.MessageSubscriber, func()) {
	return m.newSubscriberWithCleanup()
}

func TestChatServer_Reader(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	done := make(chan struct{})

	messageSubscriber := model.NewSubscriber(model.NewSubscriberID(), 1)
	cleanUp := func() {
		messageSubscriber.Close()
	}

	s := &messageServiceMock{
		sendMessage: func(ctx context.Context, text string) (*model.Message, error) {
			return &model.Message{}, nil
		},
		getMessage: func(ctx context.Context, id model.MessageID) (*model.Message, error) {
			return &model.Message{}, nil
		},
		newSubscriberWithCleanup: func() (*model.MessageSubscriber, func()) {
			//close(done)
			return messageSubscriber, cleanUp
		},
	}

	serverStream := &serverStreamingMock{
		send:       func(res *pb.Message) error { return nil },
		setHeader:  func(md metadata.MD) error { return nil },
		sendHeader: func(md metadata.MD) error { return nil },
		setTrailer: func(md metadata.MD) {},
		context:    func() context.Context { return ctx },
		sendMsg:    func(m any) error { return nil },
		recvMsg:    func(m any) error { return nil },
	}

	srv := server.NewChatServer(logger, s)
	go func() {
		err := srv.Reader(&pb.ReaderRequest{}, serverStream)
		assert.NoError(t, err)
	}()

	<-done

	messageSubscriber.Push(&model.Message{
		Text: "test",
	})

	cleanUp()

	time.Sleep(500 * time.Millisecond)

	//cancel()
}
