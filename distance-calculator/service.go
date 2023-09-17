package main

import (
	"github.com/track-tracking/types"
	"math"
)

type CalculatorServicer interface {
	CalculateDistance(data types.OBUData) (float64, error)
}

type CalculatorService struct {
	points [][]float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{
		points: make([][]float64, 0),
	}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.00
	if len(s.points) > 0 {
		prevPoints := s.points[len(s.points)-1]
		distance = calcDistance(prevPoints[0], prevPoints[1], data.Latitude, data.Longitude)
	}
	s.points = append(s.points, []float64{data.Latitude, data.Longitude})
	return distance, nil
}

func calcDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
