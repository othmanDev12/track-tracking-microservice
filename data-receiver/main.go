package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/track-tracking/types"
)

const KafkaTopic = "obudata"

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.wsHandler)
	err = http.ListenAndServe(":30000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	return &DataReceiver{
		data:     make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) dataProducer(data types.OBUData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// start a goroutine the check if the data was delivered or not
	go func() {
		for e := range dr.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	topic := KafkaTopic
	err = dr.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)
	return err
}

type DataReceiver struct {
	data     chan types.OBUData
	conn     *websocket.Conn
	producer *kafka.Producer
}

func (dr *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	con, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = con

	go dr.wsReceiverLoop()
}

func (dr *DataReceiver) wsReceiverLoop() {
	fmt.Println("New Obu Data Connected .................")
	for {
		var obuData types.OBUData
		err := dr.conn.ReadJSON(&obuData)
		if err != nil {
			log.Fatal("something wrong", err)
			continue
		}
		//fmt.Printf("received Obu Data: [%d] :: <lat %.2f  long %.2f> \n", obuData.OBUID, obuData.Latitude, obuData.Longitude)
		//dr.data <- obuData
		fmt.Println("message reciveded", obuData)
	}
}
