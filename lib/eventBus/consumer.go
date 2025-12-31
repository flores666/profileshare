package eventBus

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume(ctx context.Context, handler func([]byte) error) error
}

type consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer instance.
//
// brokers is slice of brokers addresses
func NewConsumer(brokers []string, topic, group string) Consumer {
	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: group,
		}),
	}
}

// Consume reads messages from Kafka and passes them to the handler.
//
// handler is a func that receives message in []bytes
func (c *consumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if err := handler(msg.Value); err != nil {
			// логика retry / dead letter
		}
	}
}
