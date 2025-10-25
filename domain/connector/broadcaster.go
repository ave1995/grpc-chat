package connector

import (
	"github.com/segmentio/kafka-go"
)

type Broadcaster interface {
	Broadcast(msg kafka.Message) error
}
