package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const sendInterval = 60

type OBUData struct {
	OBUID     int     `json:"obuid"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

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
	for {
		for i := 0; i < len(obuids); i++ {
			latitude, longituide := generateLatitudeAndLangitude()
			data := OBUData{
				OBUID:     obuids[i],
				Latitude:  latitude,
				Longitude: longituide,
			}
			fmt.Printf("%+v\n", data)
		}
		time.Sleep(sendInterval * time.Second)
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
