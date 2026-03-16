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
