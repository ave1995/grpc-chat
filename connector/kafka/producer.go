package kafka

import (
	"context"
	"fmt"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/segmentio/kafka-go"
)

var _ connector.Producer = (*Producer)(nil)

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config config.KafkaConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(config.Brokers...),
			Balancer: &kafka.LeastBytes{},
		}}
}

func (p *Producer) Send(ctx context.Context, topic string, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("producer send message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
