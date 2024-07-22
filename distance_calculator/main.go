package main

import (
	"log"

	"github.com/kzinthant-d3v/toll-calculator/invoice/client"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000"
)

// HTTP, GRPC, Kafka transports <- attach business logic
func main() {
	var (
		err error
		svc CalculatorServicer
	)

	svc = NewCalculatorService()
	svc = NewLoggingMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	KafkaTransport, err := NewKafkaTransport(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	KafkaTransport.Start()
}
