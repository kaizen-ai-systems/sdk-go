package kaizen

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClientRequest_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-auth-1")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"bad key"}`))
	}))
	defer server.Close()

	client := newHTTPClient(server.URL, "bad", 2*time.Second)
	_, err := client.get(context.Background(), "/v1/akuma/query")
	if err == nil {
		t.Fatal("expected error")
	}

	if _, ok := err.(*AuthError); !ok {
		t.Fatalf("expected AuthError, got %T", err)
	}
	authErr := err.(*AuthError)
	if authErr.Message != "bad key" {
		t.Fatalf("expected auth message to preserve API body, got %q", authErr.Message)
	}
	if authErr.RequestID != "req-auth-1" {
		t.Fatalf("expected request id to be surfaced, got %q", authErr.RequestID)
	}
}

func TestHTTPClientRequest_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-rate-1")
		w.Header().Set("Retry-After", "9")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"slow down"}`))
	}))
	defer server.Close()

	client := newHTTPClient(server.URL, "ok", 2*time.Second)
	_, err := client.get(context.Background(), "/v1/akuma/query")
	if err == nil {
		t.Fatal("expected error")
	}

	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
	if rateLimitErr.RetryAfter != 9 {
		t.Fatalf("expected RetryAfter 9, got %d", rateLimitErr.RetryAfter)
	}
	if rateLimitErr.Message != "slow down" {
		t.Fatalf("expected response error message, got %q", rateLimitErr.Message)
	}
	if rateLimitErr.RequestID != "req-rate-1" {
		t.Fatalf("expected request id, got %q", rateLimitErr.RequestID)
	}
}

func TestHTTPClientRequest_SendsUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "kaizen-go/"+Version {
			t.Fatalf("unexpected user agent: %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newHTTPClient(server.URL, "ok", 2*time.Second)
	if _, err := client.get(context.Background(), "/health"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
