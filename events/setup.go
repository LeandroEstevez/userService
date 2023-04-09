package events

import (
	"userMicroService/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var KafkaConfig kafka.ConfigMap

func SetUp() {
	KafkaConfig = util.LoadKafkaConfig("kafkaConfig.properties")
}
