package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type DataConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
}

func NewKafkaConsumer(topic string) (*DataConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}
	return &DataConsumer{
		consumer: c,
	}, nil

}

func (c *DataConsumer) Start() {
	logrus.Info("start consume messages .............")
	c.ReadMessageLoop()
}

func (c *DataConsumer) ReadMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error: %s", err)
			continue
		}
		obuData, err := json.Marshal(msg)
		if err != nil {
			logrus.Errorf("error while marshalling data: %s", err)
		}
		logrus.Infof("consumed data: %s", obuData)

	}
}
