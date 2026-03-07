package kaizen

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAkumaClientQueryIncludesSourceID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/akuma/query" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload AkumaQueryRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload.SourceID != "src_123" {
			t.Fatalf("expected source id to round-trip, got %q", payload.SourceID)
		}
		_, _ = w.Write([]byte(`{"sql":"select 1"}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{BaseURL: server.URL, APIKey: "test"})
	resp, err := client.Akuma.Query(context.Background(), &AkumaQueryRequest{
		Dialect:  DialectPostgres,
		Prompt:   "show one row",
		SourceID: "src_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SQL != "select 1" {
		t.Fatalf("unexpected sql: %q", resp.SQL)
	}
}

func TestAkumaClientCreateSource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/akuma/sources" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		_, _ = w.Write([]byte(`{"status":"syncing","sourceId":"src_123"}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{BaseURL: server.URL, APIKey: "test"})
	resp, err := client.Akuma.CreateSource(context.Background(), &AkumaCreateSourceRequest{
		Name:             "Warehouse",
		Dialect:          DialectPostgres,
		ConnectionString: "postgres://user:pass@db.example.com/app",
		TargetSchemas:    []string{"public"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SourceID != "src_123" {
		t.Fatalf("unexpected source id: %q", resp.SourceID)
	}
}
