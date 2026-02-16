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

// InventoryClient handles inter-service communication with the Inventory service.
type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewInventoryClient creates a new InventoryClient with the given base URL.
func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateBeerLot calls the Inventory service to create a beer lot with an initial inventory movement.
// The authToken is passed through from the original request (JWT forwarding).
func (c *InventoryClient) CreateBeerLot(ctx context.Context, authToken string, req BeerLotRequest) (*BeerLotResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling beer lot request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/beer-lots", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating beer lot request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling inventory service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("inventory service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result BeerLotResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding beer lot response: %w", err)
	}

	return &result, nil
}
