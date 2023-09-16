package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/track-tracking/types"
	"log"
	"net/http"
)

func main() {
	recv := NewDataReceiver()
	http.HandleFunc("/ws", recv.wsHandler)
	err := http.ListenAndServe(":30000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		data: make(chan types.OBUData, 128),
	}
}

type DataReceiver struct {
	data chan types.OBUData
	conn *websocket.Conn
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
		fmt.Printf("received Obu Data: [%d] :: <lat %.2f  long %.2f> \n", obuData.OBUID, obuData.Latitude, obuData.Longitude)
		dr.data <- obuData
	}
}
