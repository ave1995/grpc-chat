package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/ave1995/grpc-chat/config"
	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type producer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(config config.KafkaConfig) (store.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("create Kafka producer: %w", err)
	}

	return &producer{producer: p}, nil
}

// SendMessage implements store.Producer.
func (p *producer) SendMessage(ctx context.Context, topic string, key string, value string) error {
	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("produce message: %w", err)
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			return m.TopicPartition.Error
		} else {
			log.Printf("Delivered message to %v\n", m.TopicPartition)
			return nil
		}
	case <-ctx.Done():
		return fmt.Errorf("send message canceled or timed out: %w", ctx.Err())
	}
}

// Close implements store.Producer.
func (p *producer) Close() {
	remaining := p.producer.Flush(5000)
	if remaining > 0 {
		log.Printf("Warning: %d message(s) not delivered before close", remaining)
	}
	p.producer.Close()
}
