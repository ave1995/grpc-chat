package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/segmentio/kafka-go"
)

var _ connector.Consumer = (*Consumer)(nil)

type Consumer struct {
	logger *slog.Logger
	reader *kafka.Reader
}

func NewKafkaConsumer(logger *slog.Logger, brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) (<-chan model.Message, error) {
	msgCh := make(chan model.Message, 100)

	go func() {
		defer func() {
			c.logger.Info("[Kafka] consumer shutting down")
			if err := c.reader.Close(); err != nil {
				c.logger.Error("[Kafka] close reader error:", err)
			}
			close(msgCh)
		}()

		for {
			event, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					c.logger.Info("[Kafka] consumer context canceled, shutting down...")
					return
				}
				c.logger.Error("[Kafka] Read error: ", utils.SlogError(err))
				continue
			}

			c.logger.Info("[Kafka] received message",
				"offset", event.Offset,
				"key", string(event.Key))

			msg := model.Message{
				ID:   model.MessageID(event.Key),
				Text: string(event.Value),
			}

			msgCh <- msg
		}
	}()

	return msgCh, nil
}
