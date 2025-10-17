package main

import (
	"encoding/json"
	"log"

	"github.com/krnveersharma/browserdeck/browser"
	"github.com/krnveersharma/browserdeck/config"
	"github.com/krnveersharma/browserdeck/kafka"
	"github.com/krnveersharma/browserdeck/schemas"
)

func main() {
	kafKaConfig := config.LoadKafkaConfig()
	kafkaConsumer := kafka.NewKafkaConsumer(kafKaConfig.Brokers, kafKaConfig.Topic, kafKaConfig.GroupID)
	kafkaConsumer.Listen(func(key, value []byte) error {
		log.Printf("Consumed raw message: %s", string(value))

		var msg schemas.SessionMessage
		if err := json.Unmarshal(value, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		browserLauncher := browser.GetLauncher(msg.Browser)
		err := browserLauncher.Launch(msg.SessionID, msg.Instructions, msg.Url)
		if err != nil {
			log.Printf("error in launching browser: %v", err)
		}
		log.Printf("Launching browser: %s", msg.Browser)

		return nil
	})

}
