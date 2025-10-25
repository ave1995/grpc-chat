package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/ave1995/grpc-chat/utils"
	"github.com/segmentio/kafka-go"
)

var _ connector.Consumer = (*Consumer)(nil)

type Consumer struct {
	logger      *slog.Logger
	reader      *kafka.Reader
	broadcaster connector.Broadcaster

	wg sync.WaitGroup
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

			err = c.broadcaster.Broadcast(msg)
			if err != nil {
				c.logger.Error("[Kafka] Broadcast error: ", utils.SlogError(err))
			}
		}
	}()
}

// Stop gracefully stops the consumer and waits for the read loop to exit.
// TODO: maybe could it be done better
func (c *Consumer) Stop() error {
	c.logger.Info("[Kafka] stopping consumer...")
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("consumer stop: %w", err)
	}
	c.wg.Wait()
	c.logger.Info("[Kafka] Consumer stopped.")
	return nil
}
