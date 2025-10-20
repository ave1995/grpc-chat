package kafka

import (
	"context"

	"github.com/ave1995/grpc-chat/domain/connector"
	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) connector.Consumer {
	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

// Tohle sice nikde nepoužíváš, ale už je to takhle zvláštní. Asi nějaký příklad ne? Proč by si vracel Topic, když reader je
// nastavený na konkretní topic? Výsledek bude asi překvapivý. Key nevím jestli je taky potřeba. Vracel bych už nějaký konkretní model.
// Ten "consumer" by měl implementovat ideálně nějaký connector interface takže bude na něco určený.
func (c *consumer) ReadMessage(ctx context.Context) (topic string, key string, value string, err error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return "", "", "", err
	}
	return msg.Topic, string(msg.Key), string(msg.Value), nil
}

func (c *consumer) Close() error {
	return c.reader.Close()
}
