package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kzinthant-d3v/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	fmt.Println("Getting the invoice", id)
	invoiceReq := types.GetInvoiceRequest{
		ObuID: int32(id),
	}
	b, err := json.Marshal(&invoiceReq)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obu_id=%d", c.Endpoint, "invoice", id)
	logrus.Infof("requesting get invoice -> %s", endpoint)
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("request error")
		log.Fatal(err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var inv types.Invoice
	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	httpc := http.DefaultClient
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := httpc.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
