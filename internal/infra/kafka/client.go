package kafka

import (
	"context"
	"encoding/json"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	"friend_zone/internal/config"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(cfg config.KafkaConfig) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(cfg.Brokers...),
			Balancer:     &kafkago.Hash{},
			BatchTimeout: 50 * time.Millisecond,
		},
	}
}

func (p *Producer) PublishJSON(ctx context.Context, topic string, key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.PublishBytes(ctx, topic, key, payload)
}

func (p *Producer) PublishBytes(ctx context.Context, topic string, key string, payload []byte) error {
	return p.writer.WriteMessages(ctx, kafkago.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now().UTC(),
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func NewReader(cfg config.KafkaConfig, topic string) *kafkago.Reader {
	return kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID + "-" + topic,
		Topic:          topic,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
		StartOffset:    kafkago.FirstOffset,
	})
}
