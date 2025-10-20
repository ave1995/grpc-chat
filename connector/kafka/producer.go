package kafka

import (
	"context"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config config.KafkaConfig) connector.Producer {
	// špatný formát
	return &producer{writer: &kafka.Writer{
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
	return p.writer.WriteMessages(ctx, msg) // wrapuj errory, nebudeš vědět odkud ta chyba přišla.
}

func (p *producer) Close() error {
	return p.writer.Close()
}
