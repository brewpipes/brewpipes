package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ProcurementClient handles inter-service communication with the Procurement service.
type ProcurementClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProcurementClient creates a new ProcurementClient with the given base URL.
func NewProcurementClient(baseURL string) *ProcurementClient {
	return &ProcurementClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// PurchaseOrderLineCost holds the cost-relevant fields from a procurement PO line.
type PurchaseOrderLineCost struct {
	UUID          string `json:"uuid"`
	UnitCostCents int64  `json:"unit_cost_cents"`
	Quantity      int64  `json:"quantity"`
	QuantityUnit  string `json:"quantity_unit"`
	Currency      string `json:"currency"`
}

// batchLookupRequest is the request body for the procurement batch-lookup endpoint.
type batchLookupRequest struct {
	UUIDs []string `json:"uuids"`
}

// BatchLookupPOLines calls the Procurement service to look up multiple PO lines by UUID.
// The authToken is passed through from the original request (JWT forwarding).
func (c *ProcurementClient) BatchLookupPOLines(ctx context.Context, authToken string, uuids []string) ([]PurchaseOrderLineCost, error) {
	body, err := json.Marshal(batchLookupRequest{UUIDs: uuids})
	if err != nil {
		return nil, fmt.Errorf("marshaling batch lookup request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/purchase-order-lines/batch-lookup", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating batch lookup request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling procurement service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("procurement service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result []PurchaseOrderLineCost
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding batch lookup response: %w", err)
	}

	return result, nil
}
