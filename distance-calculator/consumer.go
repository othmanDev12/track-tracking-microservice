package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/track-tracking/types"
)

type DataConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, calcServicer CalculatorServicer) (*DataConsumer, error) {
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
		consumer:    c,
		calcService: calcServicer,
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
		var obudata types.OBUData
		err = json.Unmarshal(msg.Value, &obudata)
		if err != nil {
			logrus.Errorf("error while unmarshalling data: %s", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(obudata)
		if err != nil {
			logrus.Errorf("error while calculating distances: %s", err)
			continue
		}
		fmt.Printf("calculating distance %.2f\n", distance)
	}
}
