package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/track-tracking/types"
)

var KafkaTopic = "obudata"

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
	var (
		p   DataProducer
		err error
	)
	p, err = NewKafkaProducer()
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		data:     make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.producer.ProduceData(data)
}

type DataReceiver struct {
	data     chan types.OBUData
	conn     *websocket.Conn
	producer DataProducer
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
	//fmt.Println("New Obu Data Connected .................")
	for {
		var obuData types.OBUData
		if err := dr.conn.ReadJSON(&obuData); err != nil {
			log.Fatal("something wrong", err)
			continue
		}
		//fmt.Printf("received Obu Data: [%d] :: <lat %.2f  long %.2f> \n", obuData.OBUID, obuData.Latitude, obuData.Longitude)
		//dr.data <- obuData
		//fmt.Println("data delivered to kafka ", dr.producer.ProduceDate(obuData))
		//fmt.Println("message reciveded", obuData)
		if err := dr.produceData(obuData); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}
