package model

import (
	"time"

	pb "github.com/ave1995/grpc-chat/api/grpc/proto"
)

type Message struct {
	ID        MessageID
	Text      string
	Timestamp time.Time
}

func (m Message) ToProto() *pb.Message {
	return &pb.Message{
		Text: m.Text,
	}
}
