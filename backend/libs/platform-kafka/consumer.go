package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, m kafka.Message) error

type Consumer struct {
	reader *kafka.Reader
}

type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	})
	return &Consumer{reader: r}
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(ctx, m); err != nil {
			// дальше сам решишь: логировать/ретраить/падать
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
