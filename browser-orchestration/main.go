package main

import (
	"log"

	"github.com/krnveersharma/browserdeck/config"
	"github.com/krnveersharma/browserdeck/kafka"
)

func main() {
	kafKaConfig := config.LoadKafkaConfig()
	kafkaConsumer := kafka.NewKafkaConsumer(kafKaConfig.Brokers, kafKaConfig.Topic, kafKaConfig.GroupID)
	kafkaConsumer.Listen(func(key, value []byte) error {
		log.Printf("Consumed message: key=%s value=%s", string(key), string(value))
		return nil
	})
}
