package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/kzinthant-d3v/toll-calculator/invoice/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenPort := flag.String("port", ":8080", "port to listen")
	flag.Parse()

	aggregatorServiceAddr := flag.String("aggregator-service-addr", "http://localhost:3000", "aggregator service address")
	flag.Parse()

	var (
		client         = client.NewHTTPClient(*aggregatorServiceAddr)
		invoiceHandler = newInvoiceHandler(client)
	)

	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("gateway server started %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: client}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	invoice, err := h.client.GetInvoice(context.TODO(), 1584597223)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, invoice)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("Req")
		}(time.Now())

		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
