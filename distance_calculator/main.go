package main

import (
	"log"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const kafkaTopic = "obudata"

// HTTP, GRPC, Kafka transports <- attach business logic
func main() {
	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLoggingMiddleware(svc)

	KafkaTransport, err := NewKafkaTransport(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}

	KafkaTransport.Start()
}
