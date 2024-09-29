package newavatar

import (
	"avatar_service/internal/config"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	producer *kafka.Writer
}

func New(cfg *config.Config) (*KafkaProducer, error) {
	if cfg.Kafka.Server == "" {
		return nil, fmt.Errorf("error cfg.Kafka.Server is empty")
	}

	if cfg.Kafka.TopicForWriting == "" {
		return nil, fmt.Errorf("error cfg.Kafka.TopicForWriting is empty")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Server),
		Topic:        cfg.Kafka.TopicForWriting,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &KafkaProducer{writer}, nil
}

func (kp *KafkaProducer) Close() {
	_ = kp.producer.Close()
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, userUUID, avatarLink string) error {
	JSONMsg := struct {
		UUID string `json:"uuid"`
		Link string `json:"link"`
	}{
		UUID: userUUID,
		Link: avatarLink,
	}

	msg, err := json.Marshal(JSONMsg)
	if err != nil {
		return fmt.Errorf("error json.Marshal: %v", err)
	}

	err = kp.producer.WriteMessages(ctx, kafka.Message{Value: msg})
	if err != nil {
		return fmt.Errorf("error kp.producer.WriteMessages: %v", err)
	}

	return nil
}
