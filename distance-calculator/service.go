package main

import (
	"github.com/track-tracking/types"
	"math"
)

type CalculatorServicer interface {
	CalculateDistance(data types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoints []float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.00
	if len(s.prevPoints) > 0 {
		distance = calcDistance(s.prevPoints[0], s.prevPoints[1], data.Latitude, data.Longitude)
	}
	s.prevPoints = []float64{data.Latitude, data.Longitude}
	return distance, nil
}

func calcDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
