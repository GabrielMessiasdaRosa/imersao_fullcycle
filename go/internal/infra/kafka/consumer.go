package kafka

import (
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	ConfigMap *ckafka.ConfigMap
	Topics    []string
}

func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

func (consumer *Consumer) Consume(mdsChan chan *ckafka.Message) error {
	c, err := ckafka.NewConsumer(consumer.ConfigMap)
	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics(consumer.Topics, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Kafka consumer has been started")
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			mdsChan <- msg
		}
	}
}
