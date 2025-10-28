package connector

import (
	"github.com/ave1995/grpc-chat/domain/model"
)

type Broadcaster interface {
	Broadcast(msg *model.Message) error
}
