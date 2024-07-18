package main

import (
	"math"

	"github.com/kzinthant-d3v/toll-calculator/types"
)

type CalculatorServicer interface {
	CalculateDistance(*types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data *types.OBUData) (float64, error) {
	var distance float64

	if len(s.prevPoint) > 0 {
		distance = calculateCoordDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}

	s.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculateCoordDistance(lat1, long1, lat2, long2 float64) float64 {
	return math.Sqrt(math.Pow(lat1-lat2, 2) + math.Pow(long1-long2, 2))
}
