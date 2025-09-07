package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *KafkaConsumer) Listen(handler func(key, value []byte) error) {
	log.Println("Started kafka consumer")
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		log.Printf("Received: key=%s, value=%s, offset=%d\n",
			string(msg.Key), string(msg.Value), msg.Offset)

		if err := handler(msg.Key, msg.Value); err != nil {
			log.Println("Handler error: ", err)
		}
	}
}
