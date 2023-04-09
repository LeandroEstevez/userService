package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userMicroService/events"
	"userMicroService/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	events.SetUp()
	util.SetUpConnAndStore()
	// fmt.Println(events.KafkaConfig)
	consumer, err := kafka.NewConsumer(&events.KafkaConfig)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
	}

	err = consumer.SubscribeTopics([]string{"entries_topic"}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s", err)
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			event, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*event.TopicPartition.Topic, string(event.Key), string(event.Value))

			events.Listen(event)
		}
	}

	consumer.Close()
}
