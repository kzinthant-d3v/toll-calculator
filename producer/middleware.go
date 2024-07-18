package main

import (
	"time"

	"github.com/kzinthant-d3v/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
	next Dataproducer
}

func NewLoggingMiddleware(next Dataproducer) *LoggingMiddleware {
	return &LoggingMiddleware{
		next: next,
	}
}

func (l *LoggingMiddleware) ProduceData(data types.OBUData) error {
	start := time.Now()

	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("Producing to kafka")
	}(start)

	return l.next.ProduceData(data)
}
