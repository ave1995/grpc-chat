package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/segmentio/kafka-go"
)

var _ connector.Consumer = (*Consumer)(nil)

type Consumer struct {
	logger      *slog.Logger
	reader      *kafka.Reader
	broadcaster connector.Broadcaster
}

func NewKafkaConsumer(logger *slog.Logger, brokers []string, topic, groupID string, broadcaster connector.Broadcaster) *Consumer {
	return &Consumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		broadcaster: broadcaster,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("[Kafka] consumer shutting down")
				err := c.reader.Close()
				if err != nil {
					c.logger.Error("[Kafka] close reader error:", err)
				}
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						c.logger.Info("[Kafka] consumer context canceled, shutting down...")
						return
					}
					c.logger.Error("[Kafka] Read error: ", utils.SlogError(err))
					continue
				}

				c.logger.Info("[Kafka] received message", "offset", msg.Offset, "key", string(msg.Key))

				err = c.broadcaster.Broadcast(msg)
				if err != nil {
					c.logger.Error("[Kafka] Broadcast error: ", utils.SlogError(err))
				}
			}
		}
	}()
}
