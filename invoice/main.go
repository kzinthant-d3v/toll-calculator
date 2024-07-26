package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kzinthant-d3v/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	var (
		store          = NewMemoryStore()
		svc            = NewInvoiceAggregator(store)
		listenPorthttp = os.Getenv("AGG_GRPC_ENDPOINT")
		listenPortgrpc = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricMiddleware(svc)
	svc = NewLoggingMiddleware(svc)
	// log.Fatal((listenPortgrpc, svc))
	go func() {
		log.Fatal(makeGRPCTransport(listenPortgrpc, svc))
	}()
	makeHTTPTransport(listenPorthttp, svc)
}

func makeGRPCTransport(listenPort string, svc Aggregator) error {
	fmt.Println("Starting server on port", listenPort)
	//!!tcp should be lowercase
	ln, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer ln.Close()
	//make a new grpc server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	//register the server implementation
	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenPort string, svc Aggregator) {
	fmt.Println("Starting server on port", listenPort)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	if err := http.ListenAndServe(listenPort, nil); err != nil {
		fmt.Println("Error starting server", err)
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getting invoices")
		obuID := r.URL.Query().Get("obu_id")
		fmt.Println("obu_id", obuID)
		if obuID == "" {
			http.Error(w, "obu_id is required", http.StatusBadRequest)
			return
		}
		intObuID, err := strconv.Atoi(obuID)
		if err != nil {
			http.Error(w, "obu_id should be integer", http.StatusBadRequest)
			return
		}

		invoice, err := svc.CalculateInvoice(intObuID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, invoice)

	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(v)
}
