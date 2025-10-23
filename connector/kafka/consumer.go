package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/service/message"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/segmentio/kafka-go"
)

var _ connector.Consumer = (*Consumer)(nil)

type Consumer struct {
	logger     *slog.Logger
	reader     *kafka.Reader
	messageHub *message.MessageHub

	wg sync.WaitGroup
}

func NewKafkaConsumer(logger *slog.Logger, brokers []string, topic, groupID string, hub *message.MessageHub) *Consumer {
	return &Consumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		messageHub: hub,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.logger.Info("[Kafka] consumer started", "topic", c.reader.Config().Topic)

		for {
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

			c.messageHub.Broadcast(&model.Message{
				ID:   model.MessageID(msg.Key),
				Text: string(msg.Value),
			})
		}
	}()
}

// Stop gracefully stops the consumer and waits for the read loop to exit.
func (c *Consumer) Stop() error {
	c.logger.Info("[Kafka] stopping consumer...")
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("consumer stop: %w", err)
	}
	c.wg.Wait()
	c.logger.Info("[Kafka] Consumer stopped.")
	return nil
}
