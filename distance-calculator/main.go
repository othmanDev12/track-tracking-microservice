package main

import "log"

const kafkaTopic = "obudata"

func main() {

	var calcServicer CalculatorServicer
	calcServicer = NewCalculatorService()

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, calcServicer)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.isRunning = true
	kafkaConsumer.Start()
}
