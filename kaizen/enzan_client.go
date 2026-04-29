package kaizen

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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

// Routing gets the current smart-routing config.
func (c *EnzanClient) Routing(ctx context.Context) (*EnzanRoutingConfig, error) {
	data, err := c.http.get(ctx, "/v1/enzan/routing")
	if err != nil {
		return nil, err
	}

	var resp EnzanRoutingConfigResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan routing response: %w", err)
	}
	return &resp.Routing, nil
}

// SetRouting upserts the current smart-routing config.
func (c *EnzanClient) SetRouting(ctx context.Context, req *EnzanRoutingConfigUpsertRequest) (*EnzanRoutingConfigMutationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("routing request is required")
	}
	if req.Enabled == nil {
		return nil, fmt.Errorf("enabled is required")
	}
	data, err := c.http.post(ctx, "/v1/enzan/routing", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanRoutingConfigMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan set routing response: %w", err)
	}
	return &resp, nil
}

// RoutingSavings gets smart-routing savings for the requested window.
func (c *EnzanClient) RoutingSavings(ctx context.Context, window string) (*EnzanRoutingSavingsResponse, error) {
	if trimmed := strings.TrimSpace(window); trimmed != "" {
		switch trimmed {
		case "1h", "24h", "7d", "30d":
			window = trimmed
		default:
			return nil, fmt.Errorf("window must be one of: 1h, 24h, 7d, 30d")
		}
	}
	path := "/v1/enzan/routing/savings"
	if strings.TrimSpace(window) != "" {
		path += "?window=" + url.QueryEscape(window)
	}
	data, err := c.http.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp EnzanRoutingSavingsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan routing savings response: %w", err)
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

// TriggerPricingRefresh triggers an on-demand live-pricing refresh sweep
// (admin enzan_pricing_admin required). Fire-and-forget; poll
// ListPricingRefreshLog for completion. Returns the typed response with
// status="queued" on HTTP 202 success. The HTTP 429 (concurrency cap)
// path is surfaced as *RateLimitError; the dropped body
// ({status:"dropped",triggeredBy:...}) is preserved on err.Data so
// callers can branch on the typed shape without a separate decode.
func (c *EnzanClient) TriggerPricingRefresh(ctx context.Context) (*EnzanPricingRefreshTriggerResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/pricing/refresh", struct{}{})
	if err != nil {
		return nil, err
	}

	var resp EnzanPricingRefreshTriggerResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan pricing refresh trigger response: %w", err)
	}
	return &resp, nil
}

// ListPricingRefreshLog lists recent live-pricing refresh log entries
// (admin enzan_pricing_admin required). Server default is 50; server
// clamps limit to 1..200 and rejects non-positive values with 400. Pass
// nil to use the server default; pass a non-nil pointer to forward the
// caller's value verbatim (including 0 and negative — those will hit
// server-side validation rather than being silently dropped client-side).
func (c *EnzanClient) ListPricingRefreshLog(ctx context.Context, limit *int) ([]EnzanPricingRefreshLogEntry, error) {
	path := "/v1/enzan/pricing/refresh/log"
	if limit != nil {
		path += "?limit=" + strconv.Itoa(*limit)
	}
	data, err := c.http.get(ctx, path)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Entries []EnzanPricingRefreshLogEntry `json:"entries"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan pricing refresh log response: %w", err)
	}
	return resp.Entries, nil
}

// ListPricingProviders lists registered live-pricing sources (admin view;
// admin enzan_pricing_admin required).
func (c *EnzanClient) ListPricingProviders(ctx context.Context) ([]EnzanPricingProvider, error) {
	data, err := c.http.get(ctx, "/v1/enzan/pricing/providers")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Providers []EnzanPricingProvider `json:"providers"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan pricing providers response: %w", err)
	}
	return resp.Providers, nil
}

// UpsertPricingOffer upserts one manual (admin-authored) live-pricing offer
// (admin enzan_pricing_admin required). Exactly one of req.GPU or req.LLM
// must be set. Returns the typed response with status="upserted" on
// HTTP 201 success. The HTTP 409 stale path (a newer fetched_at row
// already exists for the same key) is surfaced as *KaizenError with
// Status=409; the stale body ({status:"stale"}) is preserved on
// err.Data so callers can branch on the typed shape without a separate
// decode.
//
// Client-side validation rejects empty string identifiers (provider,
// gpuType/model, displayName) and unset rate fields before hitting the
// wire — those are the operator mistakes the server cannot meaningfully
// surface (a forgotten float field would otherwise serialize as 0 and
// pass the server's `minimum: 0` constraint as a "free" offer). Setting
// a rate to `Float64Ptr(0)` is explicitly allowed for genuinely free
// offers; the disambiguation is the whole reason the rate fields are
// pointers.
func (c *EnzanClient) UpsertPricingOffer(ctx context.Context, req *EnzanPricingOfferUpsertRequest) (*EnzanPricingOfferUpsertResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("offer request is required")
	}
	if (req.GPU == nil) == (req.LLM == nil) {
		return nil, fmt.Errorf("exactly one of gpu or llm must be set")
	}
	// requireFiniteRate rejects nil, NaN, and Inf so wire payloads carry only
	// finite numeric values. Matches TS's Number.isFinite and Python's
	// math.isfinite. Explicit zero (free offer) is allowed.
	requireFiniteRate := func(value *float64, label string) error {
		if value == nil {
			return fmt.Errorf("%s is required (use Float64Ptr(0) for explicit free offers)", label)
		}
		if math.IsNaN(*value) || math.IsInf(*value, 0) {
			return fmt.Errorf("%s must be a finite number", label)
		}
		return nil
	}
	if req.GPU != nil {
		if strings.TrimSpace(req.GPU.Provider) == "" {
			return nil, fmt.Errorf("gpu.provider is required")
		}
		if strings.TrimSpace(req.GPU.GPUType) == "" {
			return nil, fmt.Errorf("gpu.gpuType is required")
		}
		if strings.TrimSpace(req.GPU.DisplayName) == "" {
			return nil, fmt.Errorf("gpu.displayName is required")
		}
		if err := requireFiniteRate(req.GPU.HourlyRateUSD, "gpu.hourlyRateUSD"); err != nil {
			return nil, err
		}
	}
	if req.LLM != nil {
		if strings.TrimSpace(req.LLM.Provider) == "" {
			return nil, fmt.Errorf("llm.provider is required")
		}
		if strings.TrimSpace(req.LLM.Model) == "" {
			return nil, fmt.Errorf("llm.model is required")
		}
		if strings.TrimSpace(req.LLM.DisplayName) == "" {
			return nil, fmt.Errorf("llm.displayName is required")
		}
		if err := requireFiniteRate(req.LLM.InputCostPer1KTokensUSD, "llm.inputCostPer1KTokensUSD"); err != nil {
			return nil, err
		}
		if err := requireFiniteRate(req.LLM.OutputCostPer1KTokensUSD, "llm.outputCostPer1KTokensUSD"); err != nil {
			return nil, err
		}
	}
	data, err := c.http.post(ctx, "/v1/enzan/pricing/offers", req)
	if err != nil {
		return nil, err
	}

	var resp EnzanPricingOfferUpsertResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan pricing offer upsert response: %w", err)
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

// UpdateAlert updates one configured alert.
func (c *EnzanClient) UpdateAlert(ctx context.Context, id string, alert *EnzanUpdateAlertRequest) (*EnzanAlertMutationResponse, error) {
	data, err := c.http.request(ctx, "PATCH", "/v1/enzan/alerts/"+url.PathEscape(id), alert)
	if err != nil {
		return nil, err
	}

	var resp EnzanAlertMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode enzan update alert response: %w", err)
	}
	return &resp, nil
}

// DeleteAlert deletes one configured alert.
func (c *EnzanClient) DeleteAlert(ctx context.Context, id string) error {
	_, err := c.http.request(ctx, "DELETE", "/v1/enzan/alerts/"+url.PathEscape(id), nil)
	return err
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
	case CreatableAlertCostAnomaly:
		if alert.Threshold == nil {
			return fmt.Errorf("threshold is required for alert type %s", alert.Type)
		}
		if *alert.Threshold <= 0 {
			return fmt.Errorf("threshold must be greater than 0 for alert type %s", alert.Type)
		}
		if *alert.Threshold > 10000 {
			return fmt.Errorf("threshold must be less than or equal to 10000 for alert type %s", alert.Type)
		}
		if math.Abs(math.Round(*alert.Threshold*100)-(*alert.Threshold*100)) > 1e-9 {
			return fmt.Errorf("threshold must use at most two decimal places for alert type %s", alert.Type)
		}
		if window := strings.TrimSpace(alert.Window); window == "" {
			return fmt.Errorf("window is required for alert type %s", alert.Type)
		} else if window == "1h" {
			return fmt.Errorf("window must be 24h, 7d, or 30d for alert type %s", alert.Type)
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

// UpdateAlertEndpoint updates one Enzan alert delivery webhook endpoint.
func (c *EnzanClient) UpdateAlertEndpoint(ctx context.Context, id string, req *EnzanAlertEndpointUpdateRequest) (*EnzanAlertEndpointMutationResponse, error) {
	data, err := c.http.request(ctx, "PATCH", "/v1/enzan/alerts/endpoints/"+url.PathEscape(id), req)
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
