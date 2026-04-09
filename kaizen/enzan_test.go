package kaizen

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEnzanSummaryResponseUnmarshalWithAPICosts(t *testing.T) {
	input := []byte(`{
		"window": "24h",
		"startTime": "2026-02-20T00:00:00Z",
		"endTime": "2026-02-20T23:59:59Z",
		"rows": [],
		"total": {
			"cost_usd": 12.5,
			"gpu_hours": 3.2,
			"requests": 10,
			"tokens_in": 1000,
			"tokens_out": 200
		},
		"apiCosts": {
			"totalCostUsd": 0.42,
			"promptTokens": 1000,
			"outputTokens": 200,
			"queries": 5
		}
	}`)

	var got EnzanSummaryResponse
	if err := json.Unmarshal(input, &got); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if got.APICosts == nil {
		t.Fatal("expected apiCosts to be populated")
	}
	if got.APICosts.TotalCostUSD != 0.42 {
		t.Fatalf("unexpected total cost: %v", got.APICosts.TotalCostUSD)
	}
	if got.APICosts.Queries != 5 {
		t.Fatalf("unexpected queries: %d", got.APICosts.Queries)
	}
	if got.Total.TokensIn != 1000 || got.Total.TokensOut != 200 {
		t.Fatalf("unexpected total tokens: %+v", got.Total)
	}
}

func TestEnzanSummaryResponseUnmarshalWithoutAPICosts(t *testing.T) {
	input := []byte(`{
		"window": "24h",
		"startTime": "2026-02-20T00:00:00Z",
		"endTime": "2026-02-20T23:59:59Z",
		"rows": [],
		"total": {
			"cost_usd": 0,
			"gpu_hours": 0,
			"requests": 0
		}
	}`)

	var got EnzanSummaryResponse
	if err := json.Unmarshal(input, &got); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if got.APICosts != nil {
		t.Fatalf("expected apiCosts to be nil, got %+v", got.APICosts)
	}
}

func TestEnzanModelCostResponseUnmarshal(t *testing.T) {
	input := []byte(`{
		"window": "30d",
		"startTime": "2026-03-01T00:00:00Z",
		"endTime": "2026-03-30T23:59:59Z",
		"rows": [{
			"model": "gpt-4o-mini",
			"queries": 12,
			"prompt_tokens": 1200,
			"output_tokens": 600,
			"cost_usd": 3.5,
			"percentage": 70,
			"avg_cost_per_query": 0.291666,
			"categories": [{
				"category": "simple",
				"queries": 5,
				"prompt_tokens": 300,
				"output_tokens": 120,
				"cost_usd": 0.9,
				"percentage": 25.714285,
				"avg_cost_per_query": 0.18
			}]
		}],
		"total": {
			"queries": 12,
			"prompt_tokens": 1200,
			"output_tokens": 600,
			"cost_usd": 3.5
		}
	}`)

	var got EnzanModelCostResponse
	if err := json.Unmarshal(input, &got); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if got.Total.CostUSD != 3.5 {
		t.Fatalf("unexpected total cost: %v", got.Total.CostUSD)
	}
	if len(got.Rows) != 1 || got.Rows[0].Model != "gpt-4o-mini" {
		t.Fatalf("unexpected rows: %+v", got.Rows)
	}
	if len(got.Rows[0].Categories) != 1 || got.Rows[0].Categories[0].Category != "simple" {
		t.Fatalf("unexpected categories: %+v", got.Rows[0].Categories)
	}
}

func TestEnzanPricingResponsesUnmarshal(t *testing.T) {
	listInput := []byte(`{
		"models": [{
			"provider": "openai",
			"model": "gpt-4o-mini",
			"display_name": "GPT-4o mini",
			"input_cost_per_1k_tokens_usd": 0.00015,
			"output_cost_per_1k_tokens_usd": 0.0006,
			"currency": "USD",
			"active": true
		}]
	}`)
	var listResp struct {
		Models []EnzanLLMPricing `json:"models"`
	}
	if err := json.Unmarshal(listInput, &listResp); err != nil {
		t.Fatalf("unexpected model pricing unmarshal error: %v", err)
	}
	if len(listResp.Models) != 1 || listResp.Models[0].Model != "gpt-4o-mini" {
		t.Fatalf("unexpected model pricing rows: %+v", listResp.Models)
	}

	mutationInput := []byte(`{
		"status": "upserted",
		"pricing": {
			"provider": "runpod",
			"gpu_type": "h100",
			"display_name": "H100",
			"hourly_rate_usd": 1.99,
			"currency": "USD",
			"active": true
		}
	}`)
	var mutationResp EnzanGPUPricingMutationResponse
	if err := json.Unmarshal(mutationInput, &mutationResp); err != nil {
		t.Fatalf("unexpected gpu pricing mutation unmarshal error: %v", err)
	}
	if mutationResp.Status != "upserted" || mutationResp.Pricing.GPUType != "h100" {
		t.Fatalf("unexpected gpu pricing mutation response: %+v", mutationResp)
	}
}

func TestEnzanOptimizeResponseUnmarshal(t *testing.T) {
	raw := `{
		"window":"30d",
		"startTime":"2026-03-01T00:00:00Z",
		"endTime":"2026-03-31T00:00:00Z",
		"efficiencyScore":85,
		"monthlySpend":100.50,
		"potentialSavings":15.25,
		"recommendations":[{
			"type":"model_downgrade",
			"title":"Downgrade simple queries",
			"description":"50% of queries are simple",
			"estimatedSavings":15.25,
			"confidence":0.8,
			"suggestion":"Route to cheaper model"
		}]
	}`
	var resp EnzanOptimizeResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.EfficiencyScore != 85 {
		t.Fatalf("expected score 85, got %d", resp.EfficiencyScore)
	}
	if len(resp.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(resp.Recommendations))
	}
	if resp.Recommendations[0].Type != EnzanRecModelDowngrade {
		t.Fatalf("expected type model_downgrade, got %s", resp.Recommendations[0].Type)
	}
}

func TestValidateCreateAlertRequestRequiresWindowForCostThreshold(t *testing.T) {
	threshold := 100.0
	err := validateCreateAlertRequest(&EnzanCreateAlertRequest{
		Name:      "High spend",
		Type:      CreatableAlertCostThreshold,
		Threshold: &threshold,
	})
	if err == nil {
		t.Fatal("expected validation error for missing window")
	}
	if !strings.Contains(err.Error(), "window is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCreateAlertRequestRequiresThresholdForBudgetExceeded(t *testing.T) {
	err := validateCreateAlertRequest(&EnzanCreateAlertRequest{
		Name: "Budget",
		Type: CreatableAlertBudgetExceeded,
	})
	if err == nil {
		t.Fatal("expected validation error for missing threshold")
	}
	if !strings.Contains(err.Error(), "threshold is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateCreateAlertRequestAllowsOptimizationAvailableWithoutWindow(t *testing.T) {
	err := validateCreateAlertRequest(&EnzanCreateAlertRequest{
		Name: "Optimizer",
		Type: CreatableAlertOptimizationAvailable,
	})
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateCreateAlertRequestRejectsNon24HourDailySummaryWindow(t *testing.T) {
	err := validateCreateAlertRequest(&EnzanCreateAlertRequest{
		Name:   "Daily summary",
		Type:   CreatableAlertDailySummary,
		Window: "7d",
	})
	if err == nil {
		t.Fatal("expected validation error for invalid daily_summary window")
	}
	if !strings.Contains(err.Error(), "window must be 24h") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateAlertRequestMarshalIncludesEmptyWindowAndLabels(t *testing.T) {
	window := ""
	labels := map[string]string{}
	req := EnzanUpdateAlertRequest{
		Window: &window,
		Labels: &labels,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal update alert request: %v", err)
	}
	if !strings.Contains(string(data), `"window":""`) {
		t.Fatalf("expected empty window to be preserved, got %s", string(data))
	}
	if !strings.Contains(string(data), `"labels":{}`) {
		t.Fatalf("expected empty labels object to be preserved, got %s", string(data))
	}
}

func TestUpdateAlertEndpointRequestMarshalIncludesEmptySigningSecret(t *testing.T) {
	targetURL := ""
	signingSecret := ""
	req := EnzanAlertEndpointUpdateRequest{
		TargetURL:     &targetURL,
		SigningSecret: &signingSecret,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal update alert endpoint request: %v", err)
	}
	if !strings.Contains(string(data), `"targetUrl":""`) {
		t.Fatalf("expected empty targetUrl to be preserved, got %s", string(data))
	}
	if !strings.Contains(string(data), `"signingSecret":""`) {
		t.Fatalf("expected empty signingSecret to be preserved, got %s", string(data))
	}
}
