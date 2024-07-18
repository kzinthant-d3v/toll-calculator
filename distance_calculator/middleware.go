package main

import (
	"time"

	"github.com/kzinthant-d3v/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
	next CalculatorServicer
}

func NewLoggingMiddleware(next CalculatorServicer) *LoggingMiddleware {
	return &LoggingMiddleware{
		next: next,
	}
}

func (l *LoggingMiddleware) CalculateDistance(data *types.OBUData) (dist float64, err error) {

	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": dist,
		}).Info("calculate distance")
	}(time.Now())

	dist, err = l.next.CalculateDistance(data)
	return
}
