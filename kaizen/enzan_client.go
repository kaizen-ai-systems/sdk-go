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

// CostsByModel gets model-level API cost analytics for a time window.
func (c *EnzanClient) CostsByModel(ctx context.Context, req *EnzanModelCostRequest) (*EnzanModelCostResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/costs/by-model", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanModelCostResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan costs by model response: %w", err)
	}
	return &resp, nil
}

// ListModelPricing lists configured LLM pricing entries.
func (c *EnzanClient) ListModelPricing(ctx context.Context) ([]EnzanLLMPricing, error) {
	data, err := c.http.get(ctx, "/v1/enzan/pricing/models")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Models []EnzanLLMPricing `json:"models"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan model pricing response: %w", err)
	}
	return resp.Models, nil
}

// UpsertModelPricing upserts one LLM pricing entry.
func (c *EnzanClient) UpsertModelPricing(ctx context.Context, req *EnzanLLMPricingUpsertRequest) (*EnzanLLMPricingMutationResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/pricing/models", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanLLMPricingMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan model pricing mutation response: %w", err)
	}
	return &resp, nil
}

// ListGPUPricing lists configured GPU pricing entries.
func (c *EnzanClient) ListGPUPricing(ctx context.Context) ([]EnzanGPUPricing, error) {
	data, err := c.http.get(ctx, "/v1/enzan/pricing/gpus")
	if err != nil {
		return nil, err
	}

	var resp struct {
		GPUs []EnzanGPUPricing `json:"gpus"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan GPU pricing response: %w", err)
	}
	return resp.GPUs, nil
}

// UpsertGPUPricing upserts one GPU pricing entry.
func (c *EnzanClient) UpsertGPUPricing(ctx context.Context, req *EnzanGPUPricingUpsertRequest) (*EnzanGPUPricingMutationResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/pricing/gpus", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanGPUPricingMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan GPU pricing mutation response: %w", err)
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

// Optimize generates cost optimization recommendations for a time window.
func (c *EnzanClient) Optimize(ctx context.Context, req *EnzanOptimizeRequest) (*EnzanOptimizeResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/optimize", req)
	if err != nil {
		return nil, err
	}
	var resp EnzanOptimizeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan optimize response: %w", err)
	}
	return &resp, nil
}

// Chat sends a conversational AI cost Q&A message with optional multi-turn support.
func (c *EnzanClient) Chat(ctx context.Context, req *EnzanChatRequest) (*EnzanChatResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/chat", req)
	if err != nil {
		return nil, err
	}
	var resp EnzanChatResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan chat response: %w", err)
	}
	return &resp, nil
}
