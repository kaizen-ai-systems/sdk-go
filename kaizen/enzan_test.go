package kaizen

import (
	"encoding/json"
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
