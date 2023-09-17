package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/track-tracking/types"
)

type DataProducer interface {
	ProduceData(data types.OBUData) error
}

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	// start a goroutine the check if the data was delivered or not
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					//fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					//fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &KafkaProducer{producer: p}, nil
}

func (k *KafkaProducer) ProduceData(data types.OBUData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &KafkaTopic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)
}
