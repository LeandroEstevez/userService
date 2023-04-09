package events

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var Producer *kafka.Producer

type DefaultMessage struct {
	Value string
}

func SetupProducer() {
	var err error
	Producer, err = kafka.NewProducer(&KafkaConfig)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}
}

func Produce(key string, message interface{}) {
	value, _ := json.Marshal(message)

	var err error
	Producer, err = kafka.NewProducer(&KafkaConfig)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	topic := "user_topic"

	Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	// Wait for all messages to be delivered
	Producer.Flush(15 * 1000)
	Producer.Close()
}

func ProduceToEmail(key string, message interface{}) {
	value, _ := json.Marshal(message)

	var err error
	Producer, err = kafka.NewProducer(&KafkaConfig)
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	topic := "email_topic"

	Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	// Wait for all messages to be delivered
	Producer.Flush(15 * 1000)
	Producer.Close()
}
