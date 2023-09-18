package main

import (
	"github.com/gorilla/websocket"
	"github.com/track-tracking/types"
	"log"
	"math"
	"math/rand"
	"time"
)

const socketConnection = "ws://127.0.0.1:30000/ws"

const sendInterval = time.Second * 5

func generateCordinate() float64 {
	number := float64(rand.Intn(100) + 1)
	float := rand.Float64()
	return number + float
}

func generateLatitudeAndLangitude() (float64, float64) {
	return generateCordinate(), generateCordinate()
}

func main() {
	obuids := generateOBUIds(20)

	conn, _, err := websocket.DefaultDialer.Dial(socketConnection, nil)

	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < len(obuids); i++ {
			latitude, longituide := generateLatitudeAndLangitude()
			data := types.OBUData{
				OBUID:     obuids[i],
				Latitude:  latitude,
				Longitude: longituide,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func generateOBUIds(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
