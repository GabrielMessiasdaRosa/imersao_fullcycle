package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}

func (producer *Producer) Publish(msg interface{}, key []byte, topic string) error {
	p, err := ckafka.NewProducer(producer.ConfigMap)
	if err != nil {
		return err
	}
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Key:   key,
		Value: msg.([]byte),
	}
	err = p.Produce(message, nil)
	if err != nil {
		return err
	}
	return nil
}
