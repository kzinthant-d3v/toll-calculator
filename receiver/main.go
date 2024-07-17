package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kzinthant-d3v/toll-calculator/types"
)

var kafkaTopic string = "obudata"

func main() {

	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch    chan types.OBUData
	conn     *websocket.Conn
	producer Dataproducer
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := NewKafkaProducer()
	if err != nil {
		return nil, err
	}
	return &DataReceiver{
		msgch:    make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.producer.ProduceData(data)
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("new OBU client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			fmt.Println("produce error:", err)
		}
	}

}
