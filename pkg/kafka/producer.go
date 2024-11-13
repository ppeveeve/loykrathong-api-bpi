package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer struct that holds the Kafka writer
type Producer struct {
	Writer *kafka.Writer
}

// NewProducer initializes a new Kafka producer
func NewProducer(broker string, topic string) *Producer {
	log.Printf("Initializing Kafka producer with broker: %s and topic: %s", broker, topic)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.Hash{},
		Async:    true, // Optionally, use async mode
	})

	log.Println("Kafka producer initialized")
	return &Producer{Writer: writer}
}

// PublishMessage sends a message to Kafka
func (p *Producer) PublishMessage(key string, value []byte) error {
	log.Printf("Publishing message to Kafka. Key: %s, Value: %s", key, string(value))

	err := p.Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
		return err
	}
	log.Println("Message successfully published to Kafka")
	return nil
}
