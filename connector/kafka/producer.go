package kafka

import (
	"context"
	"fmt"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/segmentio/kafka-go"
)

var _ connector.Producer = (*producer)(nil)

type producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config config.KafkaConfig) connector.Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(config.Brokers...),
			Balancer: &kafka.LeastBytes{},
		}}
}

func (p *producer) SendMessage(ctx context.Context, topic string, key string, value string) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(value),
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("producer send message: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
