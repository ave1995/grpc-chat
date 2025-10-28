package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/segmentio/kafka-go"
)

var _ connector.Producer = (*Producer)(nil)

type Producer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func NewKafkaProducer(logger *slog.Logger, config config.KafkaConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(config.Brokers...),
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}
}

func (p *Producer) Send(ctx context.Context, topic string, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}
	p.logger.Info("sending message", "topic", topic, "key", key, "value", string(value))
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("producer send message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
