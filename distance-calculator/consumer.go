package main

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/track-tracking/aggregator/client"
	"github.com/track-tracking/types"
	"time"
)

type DataConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string, calcServicer CalculatorServicer, aggClient *client.Client) (*DataConsumer, error) {
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
		aggClient:   aggClient,
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
		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUId: obudata.OBUID,
		}
		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("error while aggregating inventory: %s", err)
			continue
		}
	}
}
