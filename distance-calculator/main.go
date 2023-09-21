package main

import (
	"github.com/track-tracking/aggregator/client"
	"log"
)

const (
	kafkaTopic        = "obudata"
	aggregatorEnpoint = "http://localhost:3000/aggregate"
)

func main() {

	var calcServicer CalculatorServicer
	calcServicer = NewCalculatorService()
	calcServicer = NewLogMiddleware(calcServicer)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, calcServicer, client.NewClient(aggregatorEnpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.isRunning = true
	kafkaConsumer.Start()
}
