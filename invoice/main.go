package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/kzinthant-d3v/toll-calculator/types"
)

func main() {
	listenPort := flag.String("port", "3000", "port to listen")
	flag.Parse()
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLoggingMiddleware(svc)
	makeHTTPTransport(listenPort, svc)
}

func makeHTTPTransport(listenPort *string, svc Aggregator) {
	fmt.Println("Starting server on port", *listenPort)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	if err := http.ListenAndServe(":"+*listenPort, nil); err != nil {
		fmt.Println("Error starting server", err)
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
