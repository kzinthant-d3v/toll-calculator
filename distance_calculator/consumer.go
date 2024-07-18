package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/kzinthant-d3v/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type KafkaTransport struct {
	consumer          *kafka.Consumer
	isRunning         bool
	calculatorService CalculatorServicer
}

func NewKafkaTransport(topic string, svc CalculatorServicer) (*KafkaTransport, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		return nil, err
	}

	return &KafkaTransport{consumer: c, calculatorService: svc}, nil
}

func (c *KafkaTransport) Start() {
	logrus.Info("Starting Kafka Transport")
	c.isRunning = true
	// go c.readMessageLoop()
	c.readMessageLoop()
}

func (c *KafkaTransport) readMessageLoop() {
	for c.isRunning {
		// msg, err := c.consumer.ReadMessage(time.Second)
		msg, err := c.consumer.ReadMessage(-1)
		fmt.Println(msg)
		if err != nil {
			logrus.Errorf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		var data *types.OBUData
		//in production, change this to protobuf
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON unmarshal error: %v\n", err)
			continue
		}
		distance, err := c.calculatorService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Error calculating distance: %v\n", err)
			continue
		}
		_ = distance
	}
}
