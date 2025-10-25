package server_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
	"github.com/ave1995/grpc-chat/api/grpc/server"
	"github.com/ave1995/grpc-chat/domain/model"
	"google.golang.org/grpc/metadata"
)

// --- Fake gRPC stream implementation ---
var _ pb.ChatService_ReaderServer = (*fakeReaderServer)(nil)

type fakeReaderServer struct {
	ctx  context.Context
	sent []*pb.Message
}

func (f *fakeReaderServer) Context() context.Context { return f.ctx }
func (f *fakeReaderServer) Send(msg *pb.Message) error {
	f.sent = append(f.sent, msg)
	return nil
}

// no-op implementations for required methods
func (f *fakeReaderServer) RecvMsg(m any) error             { return nil }
func (f *fakeReaderServer) SendHeader(md metadata.MD) error { return nil }
func (f *fakeReaderServer) SetHeader(md metadata.MD) error  { return nil }
func (f *fakeReaderServer) SetTrailer(md metadata.MD)       {}
func (f *fakeReaderServer) SendMsg(m any) error             { return nil }

type fakeMessageService struct {
	subscriber *model.MessageSubscriber
}

func (f *fakeMessageService) Send(ctx context.Context, text string) (*model.Message, error) {
	return nil, nil
}

func (f *fakeMessageService) Fetch(ctx context.Context, id model.MessageID) (*model.Message, error) {
	return nil, nil
}

func (f *fakeMessageService) NewSubscriberWithCleanup() (*model.MessageSubscriber, func()) {
	return f.subscriber, func() {}
}

// --- Helper: slog logger for tests ---

func newLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	})
	return slog.New(handler)
}

// --- TEST 1: Closed channel case ---

func TestReader_ChannelClosed(t *testing.T) {
	// Create real subscriber and close its channel
	sub := model.NewSubscriber(model.NewSubscriberID(), 1)
	sub.Close() // simulate hub closed

	msgSvc := &fakeMessageService{subscriber: sub}
	logger := newLogger()

	srv := server.NewChatServer(logger, msgSvc)

	stream := &fakeReaderServer{ctx: context.Background()}

	err := srv.Reader(&pb.ReaderRequest{}, stream)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stream.sent) != 0 {
		t.Fatalf("expected 0 messages sent, got %d", len(stream.sent))
	}
}

// --- TEST 2: Send one message, then close channel ---

func TestReader_SendThenClose(t *testing.T) {
	sub := model.NewSubscriber(model.NewSubscriberID(), 2)
	sub.Push(&model.Message{Text: "hello"})
	sub.Close()

	msgSvc := &fakeMessageService{subscriber: sub}
	logger := newLogger()
	srv := server.NewChatServer(logger, msgSvc)

	stream := &fakeReaderServer{ctx: context.Background()}

	err := srv.Reader(&pb.ReaderRequest{}, stream)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stream.sent) != 1 {
		t.Fatalf("expected 1 message sent, got %d", len(stream.sent))
	}
	if stream.sent[0].Text != "hello" {
		t.Fatalf("expected message 'hello', got %q", stream.sent[0].Text)
	}
}
