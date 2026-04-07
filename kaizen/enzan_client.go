package kaizen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

// CreateAlert creates an alert using the currently live Phase 6a evaluator subset.
func (c *EnzanClient) CreateAlert(ctx context.Context, alert *EnzanCreateAlertRequest) (*StatusWithIDResponse, error) {
	if err := validateCreateAlertRequest(alert); err != nil {
		return nil, err
	}
	data, err := c.http.post(ctx, "/v1/enzan/alerts", alert)
	if err != nil {
		return nil, err
	}

	var resp StatusWithIDResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan create alert response: %w", err)
	}
	return &resp, nil
}

func validateCreateAlertRequest(alert *EnzanCreateAlertRequest) error {
	if alert == nil {
		return fmt.Errorf("alert is required")
	}
	switch alert.Type {
	case CreatableAlertCostThreshold:
		if alert.Threshold == nil {
			return fmt.Errorf("threshold is required for alert type %s", alert.Type)
		}
		if strings.TrimSpace(alert.Window) == "" {
			return fmt.Errorf("window is required for alert type %s", alert.Type)
		}
	case CreatableAlertBudgetExceeded:
		if alert.Threshold == nil {
			return fmt.Errorf("threshold is required for alert type %s", alert.Type)
		}
	case CreatableAlertDailySummary:
		if window := strings.TrimSpace(alert.Window); window != "" && window != "24h" {
			return fmt.Errorf("window must be 24h for alert type %s", alert.Type)
		}
	}
	return nil
}

// ListAlertEndpoints lists configured Enzan alert delivery endpoints.
func (c *EnzanClient) ListAlertEndpoints(ctx context.Context) ([]EnzanAlertEndpoint, error) {
	data, err := c.http.get(ctx, "/v1/enzan/alerts/endpoints")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Endpoints []EnzanAlertEndpoint `json:"endpoints"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan alert endpoints response: %w", err)
	}
	return resp.Endpoints, nil
}

// CreateAlertEndpoint creates one Enzan alert delivery webhook endpoint.
func (c *EnzanClient) CreateAlertEndpoint(ctx context.Context, req *EnzanAlertEndpointCreateRequest) (*EnzanAlertEndpointMutationResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/alerts/endpoints", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanAlertEndpointMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan alert endpoint mutation response: %w", err)
	}
	return &resp, nil
}

// ListAlertEvents lists recent Enzan alert events.
func (c *EnzanClient) ListAlertEvents(ctx context.Context, limit int) ([]EnzanAlertEvent, error) {
	path := "/v1/enzan/alerts/events"
	if limit > 0 {
		path += "?limit=" + strconv.Itoa(limit)
	}
	data, err := c.http.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Events []EnzanAlertEvent `json:"events"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan alert events response: %w", err)
	}
	return resp.Events, nil
}

// ListAlertDeliveries lists recent Enzan alert deliveries.
func (c *EnzanClient) ListAlertDeliveries(ctx context.Context, limit int) ([]EnzanAlertDelivery, error) {
	path := "/v1/enzan/alerts/deliveries"
	if limit > 0 {
		path += "?limit=" + strconv.Itoa(limit)
	}
	data, err := c.http.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Deliveries []EnzanAlertDelivery `json:"deliveries"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan alert deliveries response: %w", err)
	}
	return resp.Deliveries, nil
}

// DeleteAlertEndpoint deletes one Enzan alert delivery webhook endpoint.
func (c *EnzanClient) DeleteAlertEndpoint(ctx context.Context, id string) error {
	_, err := c.http.request(ctx, "DELETE", "/v1/enzan/alerts/endpoints/"+url.PathEscape(id), nil)
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
