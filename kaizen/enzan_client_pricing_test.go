package kaizen

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEnzanLivePricingAdminClientMethods(t *testing.T) {
	var capturedOfferBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/enzan/pricing/refresh":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"status":"queued","triggeredBy":"33333333-3333-3333-3333-333333333333"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/enzan/pricing/refresh/log":
			if got := r.URL.Query().Get("limit"); got != "5" {
				t.Fatalf("unexpected limit query: got %q want %q", got, "5")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"entries":[{"id":"11111111-1111-1111-1111-111111111111","kind":"on_demand","status":"success","rowsUpserted":0,"rowsSkipped":0,"durationMs":64,"startedAt":"2026-04-28T13:56:13.416941Z","finishedAt":"2026-04-28T13:56:13.483386Z","sourceId":"22222222-2222-2222-2222-222222222222","sourceName":"manual","triggeredBy":"33333333-3333-3333-3333-333333333333"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/enzan/pricing/providers":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"providers":[{"id":"44444444-4444-4444-4444-444444444444","name":"manual","kind":"manual","enabled":true,"refreshIntervalHours":24,"hasAdapter":true}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/enzan/pricing/offers":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			if err := json.Unmarshal(body, &capturedOfferBody); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"status":"upserted","gpu":{"id":"55555555-5555-5555-5555-555555555555","provider":"manual-smoke","gpuType":"h100-80gb","displayName":"Smoke H100","deploymentClass":"on_demand","clusterSizeMin":1,"interconnectClass":"unknown","trainingReady":false,"hourlyRateUSD":2.99,"currency":"USD","sourceType":"admin","trustStatus":"verified","fetchedAt":"2026-04-28T13:57:38Z","firstSeenAt":"2026-04-28T13:57:38Z","lastSeenAt":"2026-04-28T13:57:38Z","active":true}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})

	ctx := context.Background()

	triggered, err := client.Enzan.TriggerPricingRefresh(ctx)
	if err != nil {
		t.Fatalf("TriggerPricingRefresh() error = %v", err)
	}
	if triggered.Status != "queued" || triggered.TriggeredBy != "33333333-3333-3333-3333-333333333333" {
		t.Fatalf("unexpected refresh trigger response: %+v", triggered)
	}

	limitFive := 5
	log, err := client.Enzan.ListPricingRefreshLog(ctx, &limitFive)
	if err != nil {
		t.Fatalf("ListPricingRefreshLog() error = %v", err)
	}
	if len(log) != 1 || log[0].Kind != "on_demand" || log[0].Status != "success" {
		t.Fatalf("unexpected refresh log response: %+v", log)
	}
	if log[0].SourceName == nil || *log[0].SourceName != "manual" {
		t.Fatalf("expected sourceName=manual, got %+v", log[0].SourceName)
	}

	providers, err := client.Enzan.ListPricingProviders(ctx)
	if err != nil {
		t.Fatalf("ListPricingProviders() error = %v", err)
	}
	if len(providers) != 1 || !providers[0].HasAdapter || providers[0].Kind != "manual" {
		t.Fatalf("unexpected providers response: %+v", providers)
	}

	upsertResp, err := client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{
		GPU: &EnzanGPUOfferUpsertPayload{
			Provider:        "manual-smoke",
			GPUType:         "h100-80gb",
			DisplayName:     "Smoke H100",
			DeploymentClass: "on_demand",
			HourlyRateUSD:   Float64Ptr(2.99),
			Currency:        "USD",
		},
	})
	if err != nil {
		t.Fatalf("UpsertPricingOffer() error = %v", err)
	}
	if upsertResp.Status != "upserted" || upsertResp.GPU == nil || upsertResp.GPU.SourceType != "admin" {
		t.Fatalf("unexpected upsert response: %+v", upsertResp)
	}
	if upsertResp.GPU.DeploymentClass != "on_demand" {
		t.Fatalf("unexpected deployment class: %q", upsertResp.GPU.DeploymentClass)
	}
	gpuPayload, ok := capturedOfferBody["gpu"].(map[string]any)
	if !ok {
		t.Fatalf("expected gpu payload in request body, got: %+v", capturedOfferBody)
	}
	if gpuPayload["deploymentClass"] != "on_demand" || gpuPayload["hourlyRateUSD"] != 2.99 {
		t.Fatalf("unexpected gpu payload: %+v", gpuPayload)
	}
}

func TestEnzanUpsertPricingOfferRejectsUnsetRateFields(t *testing.T) {
	// Codex pass 7 finding: Go's float64 zero value would silently submit
	// $0/hr offers from caller typos. Rate fields are *float64 and validated
	// non-nil. Float64Ptr(0) remains valid for genuinely-free offers.
	client := NewClient(&ClientConfig{
		BaseURL: "https://example.invalid",
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	ctx := context.Background()

	_, err := client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{
		GPU: &EnzanGPUOfferUpsertPayload{
			Provider:    "p",
			GPUType:     "g",
			DisplayName: "d",
			// HourlyRateUSD intentionally omitted — must be rejected.
		},
	})
	if err == nil || !strings.Contains(err.Error(), "gpu.hourlyRateUSD is required") {
		t.Fatalf("expected hourlyRateUSD-required error, got %v", err)
	}

	_, err = client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{
		LLM: &EnzanLLMOfferUpsertPayload{
			Provider:    "p",
			Model:       "m",
			DisplayName: "d",
			// Cost fields intentionally omitted.
		},
	})
	if err == nil || !strings.Contains(err.Error(), "llm.inputCostPer1KTokensUSD is required") {
		t.Fatalf("expected inputCostPer1KTokensUSD-required error, got %v", err)
	}
}

func TestEnzanUpsertPricingOfferRejectsNonFiniteRate(t *testing.T) {
	// NaN/Inf must be rejected client-side (matches TS Number.isFinite and
	// Python math.isfinite) so they never reach the wire as malformed JSON.
	client := NewClient(&ClientConfig{
		BaseURL: "https://example.invalid",
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	ctx := context.Background()
	nan := math.NaN()
	inf := math.Inf(1)

	_, err := client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{
		GPU: &EnzanGPUOfferUpsertPayload{
			Provider:      "p",
			GPUType:       "g",
			DisplayName:   "d",
			HourlyRateUSD: &nan,
		},
	})
	if err == nil || !strings.Contains(err.Error(), "must be a finite number") {
		t.Fatalf("expected NaN to be rejected, got err=%v", err)
	}

	_, err = client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{
		LLM: &EnzanLLMOfferUpsertPayload{
			Provider:                 "p",
			Model:                    "m",
			DisplayName:              "d",
			InputCostPer1KTokensUSD:  &inf,
			OutputCostPer1KTokensUSD: Float64Ptr(0),
		},
	})
	if err == nil || !strings.Contains(err.Error(), "must be a finite number") {
		t.Fatalf("expected Inf to be rejected, got err=%v", err)
	}
}

func TestEnzanUpsertPricingOfferAllowsExplicitZeroFreeOffer(t *testing.T) {
	// Float64Ptr(0) is valid — the spec accepts minimum: 0 and free offers
	// are a real concept. The pointer disambiguates this from the typo case.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"hourlyRateUSD":0`) {
			t.Fatalf("expected explicit hourlyRateUSD=0 in body, got %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"upserted","gpu":{"id":"x","provider":"free","gpuType":"g","displayName":"d","deploymentClass":"on_demand","clusterSizeMin":1,"interconnectClass":"unknown","trainingReady":false,"hourlyRateUSD":0,"currency":"USD","sourceType":"admin","trustStatus":"verified","fetchedAt":"2026-04-28T13:00:00Z","firstSeenAt":"2026-04-28T13:00:00Z","lastSeenAt":"2026-04-28T13:00:00Z","active":true}}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	resp, err := client.Enzan.UpsertPricingOffer(context.Background(), &EnzanPricingOfferUpsertRequest{
		GPU: &EnzanGPUOfferUpsertPayload{
			Provider:      "free",
			GPUType:       "g",
			DisplayName:   "d",
			HourlyRateUSD: Float64Ptr(0),
		},
	})
	if err != nil {
		t.Fatalf("explicit zero offer should succeed, got error %v", err)
	}
	if resp.GPU == nil || resp.GPU.HourlyRateUSD != 0 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestEnzanUpsertPricingOfferLLMHappyPath(t *testing.T) {
	var capturedBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &capturedBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"upserted","llm":{"id":"66666666-6666-6666-6666-666666666666","provider":"manual-smoke","model":"smoke-llm","displayName":"Smoke LLM","inputCostPer1KTokensUSD":0.001,"outputCostPer1KTokensUSD":0.002,"currency":"USD","sourceType":"admin","trustStatus":"verified","fetchedAt":"2026-04-28T13:00:00Z","firstSeenAt":"2026-04-28T13:00:00Z","lastSeenAt":"2026-04-28T13:00:00Z","active":true}}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	resp, err := client.Enzan.UpsertPricingOffer(context.Background(), &EnzanPricingOfferUpsertRequest{
		LLM: &EnzanLLMOfferUpsertPayload{
			Provider:                 "manual-smoke",
			Model:                    "smoke-llm",
			DisplayName:              "Smoke LLM",
			InputCostPer1KTokensUSD:  Float64Ptr(0.001),
			OutputCostPer1KTokensUSD: Float64Ptr(0.002),
			Currency:                 "USD",
		},
	})
	if err != nil {
		t.Fatalf("UpsertPricingOffer(LLM) error = %v", err)
	}
	if resp.LLM == nil || resp.LLM.Model != "smoke-llm" || resp.GPU != nil {
		t.Fatalf("unexpected LLM response: %+v", resp)
	}
	if resp.LLM.InputCostPer1KTokensUSD != 0.001 || resp.LLM.OutputCostPer1KTokensUSD != 0.002 {
		t.Fatalf("LLM rate fields not preserved: %+v", resp.LLM)
	}
	llmPayload, ok := capturedBody["llm"].(map[string]any)
	if !ok {
		t.Fatalf("expected llm payload in request, got: %+v", capturedBody)
	}
	if llmPayload["model"] != "smoke-llm" || llmPayload["inputCostPer1KTokensUSD"] != 0.001 {
		t.Fatalf("LLM request payload not formed correctly: %+v", llmPayload)
	}
	if _, hasGPU := capturedBody["gpu"]; hasGPU {
		t.Fatalf("LLM request must not include a gpu key: %+v", capturedBody)
	}
}

func TestEnzanUpsertPricingOfferRejectsBothOrNeither(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseURL: "https://example.invalid",
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	ctx := context.Background()

	_, err := client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{})
	if err == nil || !strings.Contains(err.Error(), "exactly one of gpu or llm") {
		t.Fatalf("UpsertPricingOffer() error = %v, want 'exactly one of gpu or llm'", err)
	}

	_, err = client.Enzan.UpsertPricingOffer(ctx, &EnzanPricingOfferUpsertRequest{
		GPU: &EnzanGPUOfferUpsertPayload{Provider: "p", GPUType: "g", DisplayName: "d", HourlyRateUSD: Float64Ptr(1)},
		LLM: &EnzanLLMOfferUpsertPayload{Provider: "p", Model: "m", DisplayName: "d", InputCostPer1KTokensUSD: Float64Ptr(0), OutputCostPer1KTokensUSD: Float64Ptr(0)},
	})
	if err == nil || !strings.Contains(err.Error(), "exactly one of gpu or llm") {
		t.Fatalf("UpsertPricingOffer() error = %v, want 'exactly one of gpu or llm'", err)
	}

	_, err = client.Enzan.UpsertPricingOffer(ctx, nil)
	if err == nil || !strings.Contains(err.Error(), "offer request is required") {
		t.Fatalf("UpsertPricingOffer(nil) error = %v, want 'offer request is required'", err)
	}
}

func TestEnzanListPricingRefreshLogPassesLimitThroughForServerClamping(t *testing.T) {
	// The client must NOT clamp client-side. Server is the authority on the
	// 1..200 range and 400-rejects non-positive values; passing values
	// through unmodified keeps server-side validation observable.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "500" {
			t.Fatalf("expected limit=500 to be forwarded as-is, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"entries":[]}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	limit500 := 500
	if _, err := client.Enzan.ListPricingRefreshLog(context.Background(), &limit500); err != nil {
		t.Fatalf("ListPricingRefreshLog() error = %v", err)
	}
}

func TestEnzanListPricingRefreshLogForwardsNegativeLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "-1" {
			t.Fatalf("expected limit=-1 to be forwarded as-is, got %q", got)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"limit must be a positive integer"}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	negOne := -1
	if _, err := client.Enzan.ListPricingRefreshLog(context.Background(), &negOne); err == nil {
		t.Fatalf("expected 400 error from server, got nil")
	}
}

func TestEnzanListPricingRefreshLogForwardsZeroLimitToTriggerServerValidation(t *testing.T) {
	// Codex-flagged: prior behavior dropped the limit query when limit<=0,
	// which prevented callers from observing the server's "limit must be a
	// positive integer" 400 path. With *int we forward whatever the caller
	// passed.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "0" {
			t.Fatalf("expected limit=0 to be forwarded as-is, got %q", got)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"limit must be a positive integer"}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	zero := 0
	_, err := client.Enzan.ListPricingRefreshLog(context.Background(), &zero)
	if err == nil {
		t.Fatalf("expected 400 error from server, got nil")
	}
}

func TestEnzanListPricingRefreshLogOmitsParamWhenNil(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("limit") {
			t.Fatalf("expected no limit param when caller passes nil, got query=%q", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"entries":[]}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	if _, err := client.Enzan.ListPricingRefreshLog(context.Background(), nil); err != nil {
		t.Fatalf("ListPricingRefreshLog(nil) error = %v", err)
	}
}

func TestEnzanTriggerPricingRefreshPreservesDroppedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"status":"dropped","triggeredBy":"33333333-3333-3333-3333-333333333333"}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	// 429 surfaces as *RateLimitError. The dropped-path body
	// ({status:"dropped",triggeredBy:"..."}) is preserved on err.Data so
	// callers can read the typed body fields without a separate decode.
	_, err := client.Enzan.TriggerPricingRefresh(context.Background())
	if err == nil {
		t.Fatalf("TriggerPricingRefresh() expected 429 error, got nil")
	}
	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("TriggerPricingRefresh() error = %v, want *RateLimitError", err)
	}
	if rateLimitErr.Status != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", rateLimitErr.Status)
	}
	if got, _ := rateLimitErr.Data["status"].(string); got != "dropped" {
		t.Fatalf("expected err.Data.status=\"dropped\", got %q (data=%+v)", got, rateLimitErr.Data)
	}
	if got, _ := rateLimitErr.Data["triggeredBy"].(string); got != "33333333-3333-3333-3333-333333333333" {
		t.Fatalf("expected err.Data.triggeredBy preserved, got %q", got)
	}
}

func TestEnzanRefreshLogHandlesNullableSourceFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"entries":[{"id":"11111111-1111-1111-1111-111111111111","kind":"scheduled","status":"failed","rowsUpserted":0,"rowsSkipped":0,"startedAt":"2026-04-28T13:00:00Z","error":"source removed mid-sweep"}]}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	log, err := client.Enzan.ListPricingRefreshLog(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListPricingRefreshLog() error = %v", err)
	}
	if len(log) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log))
	}
	entry := log[0]
	if entry.SourceID != nil || entry.SourceName != nil || entry.TriggeredBy != nil || entry.DurationMs != nil || entry.FinishedAt != nil {
		t.Fatalf("expected nullable fields to be nil, got %+v", entry)
	}
	if entry.Error == nil || *entry.Error != "source removed mid-sweep" {
		t.Fatalf("expected error string, got %+v", entry.Error)
	}
}

func TestEnzanUpsertPricingOfferReturnsStaleWithoutPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"status":"stale"}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	// 409 is surfaced as a typed *KaizenError with Status=409; callers
	// branch on Status to recognise the stale path.
	_, err := client.Enzan.UpsertPricingOffer(context.Background(), &EnzanPricingOfferUpsertRequest{
		GPU: &EnzanGPUOfferUpsertPayload{
			Provider:        "manual-smoke",
			GPUType:         "h100-80gb",
			DisplayName:     "Smoke H100",
			DeploymentClass: "on_demand",
			HourlyRateUSD:   Float64Ptr(2.99),
		},
	})
	if err == nil {
		t.Fatalf("UpsertPricingOffer() expected 409 error, got nil")
	}
	var apiErr *KaizenError
	if !errors.As(err, &apiErr) {
		t.Fatalf("UpsertPricingOffer() error = %v, want *KaizenError", err)
	}
	if apiErr.Status != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", apiErr.Status)
	}
	if got, _ := apiErr.Data["status"].(string); got != "stale" {
		t.Fatalf("expected err.Data.status=\"stale\", got %q (data=%+v)", got, apiErr.Data)
	}
}

func TestEnzanPricingProvidersHandlesOptionalFreshnessFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providers":[{"id":"44444444-4444-4444-4444-444444444444","name":"aws","kind":"api","enabled":true,"refreshIntervalHours":24,"hasAdapter":false}]}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})
	providers, err := client.Enzan.ListPricingProviders(context.Background())
	if err != nil {
		t.Fatalf("ListPricingProviders() error = %v", err)
	}
	if len(providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(providers))
	}
	p := providers[0]
	if p.HasAdapter {
		t.Fatalf("expected hasAdapter=false for aws in 8.2-public")
	}
	if p.LastSuccessAt != nil || p.LastFailureAt != nil || p.LastError != nil {
		t.Fatalf("expected freshness fields to be nil, got %+v", p)
	}
}
