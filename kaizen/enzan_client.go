package kaizen

import (
	"context"
	"encoding/json"
	"fmt"
)

// EnzanClient is the client for Enzan API.
type EnzanClient struct {
	http *httpClient
}

// Summary gets GPU cost summary for a time window.
func (c *EnzanClient) Summary(ctx context.Context, req *EnzanSummaryRequest) (*EnzanSummaryResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/summary", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanSummaryResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan summary response: %w", err)
	}
	return &resp, nil
}

// Burn gets current burn rate.
func (c *EnzanClient) Burn(ctx context.Context) (*EnzanBurnResponse, error) {
	data, err := c.http.get(ctx, "/v1/enzan/burn")
	if err != nil {
		return nil, err
	}

	var resp EnzanBurnResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan burn response: %w", err)
	}
	return &resp, nil
}

// ListResources lists registered GPU resources.
func (c *EnzanClient) ListResources(ctx context.Context) ([]EnzanResource, error) {
	data, err := c.http.get(ctx, "/v1/enzan/resources")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Resources []EnzanResource `json:"resources"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan resources response: %w", err)
	}
	return resp.Resources, nil
}

// RegisterResource registers a GPU resource.
func (c *EnzanClient) RegisterResource(ctx context.Context, resource *EnzanResource) error {
	_, err := c.http.post(ctx, "/v1/enzan/resources", resource)
	return err
}

// ListAlerts lists configured alerts.
func (c *EnzanClient) ListAlerts(ctx context.Context) ([]EnzanAlert, error) {
	data, err := c.http.get(ctx, "/v1/enzan/alerts")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Alerts []EnzanAlert `json:"alerts"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan alerts response: %w", err)
	}
	return resp.Alerts, nil
}

// CreateAlert creates an alert.
func (c *EnzanClient) CreateAlert(ctx context.Context, alert *EnzanAlert) error {
	_, err := c.http.post(ctx, "/v1/enzan/alerts", alert)
	return err
}
