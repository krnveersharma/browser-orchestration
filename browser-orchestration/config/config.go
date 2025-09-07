package config

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func LoadKafkaConfig() *KafkaConfig {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersEnv, ",")

	return &KafkaConfig{
		Brokers: brokers,
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: os.Getenv("KAFKA_GROUPID"),
	}
}
