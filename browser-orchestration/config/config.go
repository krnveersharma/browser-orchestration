package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func LoadKafkaConfig() *KafkaConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersEnv, ",")
	return &KafkaConfig{
		Brokers: brokers,
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: os.Getenv("KAFKA_GROUPID"),
	}
}
