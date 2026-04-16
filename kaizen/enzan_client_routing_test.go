package kaizen

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEnzanRoutingClientMethods(t *testing.T) {
	var capturedRoutingBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/enzan/routing":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"routing":{"enabled":true,"provider":"openai","default_model":"gpt-4.1","simple_model":"gpt-4o-mini","updated_at":"2026-04-16T12:00:00Z"}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/enzan/routing":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			if err := json.Unmarshal(body, &capturedRoutingBody); err != nil {
				t.Fatalf("decode request body: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"upserted","routing":{"enabled":true,"provider":"openai","default_model":"gpt-4.1","simple_model":"gpt-4o-mini"}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/enzan/routing/savings":
			if got := r.URL.Query().Get("window"); got != "7d" {
				t.Fatalf("unexpected window query: got %q want %q", got, "7d")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"window":"7d","start_time":"2026-04-09T00:00:00Z","end_time":"2026-04-16T00:00:00Z","provider":"openai","default_model":"gpt-4.1","total_queries":12,"routed_queries":8,"actual_cost_usd":1.2,"counterfactual_cost_usd":2.4,"estimated_savings_usd":1.2,"breakdown":[{"prompt_category":"simple","original_model":"gpt-4.1","routed_model":"gpt-4o-mini","queries":8,"actual_cost_usd":1.2,"counterfactual_cost_usd":2.4,"estimated_savings_usd":1.2}]}`))
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
	routing, err := client.Enzan.Routing(ctx)
	if err != nil {
		t.Fatalf("Routing() error = %v", err)
	}
	if !routing.Enabled || routing.SimpleModel != "gpt-4o-mini" {
		t.Fatalf("unexpected routing response: %+v", routing)
	}

	reqModel := "gpt-4o-mini"
	enabled := true
	mutation, err := client.Enzan.SetRouting(ctx, &EnzanRoutingConfigUpsertRequest{
		Enabled:     &enabled,
		SimpleModel: &reqModel,
	})
	if err != nil {
		t.Fatalf("SetRouting() error = %v", err)
	}
	if mutation.Status != "upserted" || mutation.Routing.DefaultModel != "gpt-4.1" {
		t.Fatalf("unexpected set routing response: %+v", mutation)
	}
	if capturedRoutingBody["enabled"] != true || capturedRoutingBody["simple_model"] != "gpt-4o-mini" {
		t.Fatalf("unexpected SetRouting request body: %+v", capturedRoutingBody)
	}

	savings, err := client.Enzan.RoutingSavings(ctx, "7d")
	if err != nil {
		t.Fatalf("RoutingSavings() error = %v", err)
	}
	if savings.RoutedQueries != 8 || len(savings.Breakdown) != 1 {
		t.Fatalf("unexpected routing savings response: %+v", savings)
	}
}

func TestEnzanSetRoutingRequiresExplicitEnabled(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseURL: "https://example.invalid",
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})

	_, err := client.Enzan.SetRouting(context.Background(), &EnzanRoutingConfigUpsertRequest{})
	if err == nil || !strings.Contains(err.Error(), "enabled is required") {
		t.Fatalf("SetRouting() error = %v, want enabled is required", err)
	}
}

func TestEnzanRoutingSavingsRejectsInvalidWindow(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseURL: "https://example.invalid",
		APIKey:  "test-key",
		Timeout: 2 * time.Second,
	})

	_, err := client.Enzan.RoutingSavings(context.Background(), "14d")
	if err == nil || !strings.Contains(err.Error(), "window must be one of") {
		t.Fatalf("RoutingSavings() error = %v, want window validation error", err)
	}
}
