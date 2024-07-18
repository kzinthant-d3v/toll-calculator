package main

import (
	"time"

	"github.com/kzinthant-d3v/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
	next Aggregator
}

func NewLoggingMiddleware(next Aggregator) *LoggingMiddleware {
	return &LoggingMiddleware{
		next: next,
	}
}

func (l *LoggingMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("calculate distance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return
}
