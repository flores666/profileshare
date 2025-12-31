package eventBus

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(ctx context.Context, topic string, data any) error
}

type producer struct {
	writer *kafka.Writer
}

// NewProducer creates new event bus producer.
//
// brokers is slice of brokers addresses
func NewProducer(brokers []string) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			RequiredAcks:           kafka.RequireAll,
			AllowAutoTopicCreation: true,
		},
	}
}

// Produce writes message to topic
func (p producer) Produce(ctx context.Context, topic string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: bytes,
	})
}
