package kaizen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type httpClient struct {
	mu         sync.RWMutex
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func newHTTPClient(baseURL, apiKey string, timeout time.Duration) *httpClient {
	return &httpClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *httpClient) request(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	baseURL, apiKey := c.snapshotConfig()
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("kaizen-go/%s", Version))
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	requestID := resp.Header.Get("X-Request-ID")

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, &AuthError{
			KaizenError{
				Message:   parseAPIErrorMessage(resp.StatusCode, respBody),
				Status:    http.StatusUnauthorized,
				Code:      "AUTH_ERROR",
				RequestID: requestID,
			},
		}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 0
		if raw := resp.Header.Get("Retry-After"); raw != "" {
			if parsed, parseErr := strconv.Atoi(raw); parseErr == nil {
				retryAfter = parsed
			}
		}
		return nil, &RateLimitError{
			KaizenError: KaizenError{
				Message:   parseAPIErrorMessage(resp.StatusCode, respBody),
				Status:    http.StatusTooManyRequests,
				Code:      "RATE_LIMIT",
				RequestID: requestID,
			},
			RetryAfter: retryAfter,
		}
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &KaizenError{
			Message:   parseAPIErrorMessage(resp.StatusCode, respBody),
			Status:    resp.StatusCode,
			RequestID: requestID,
			Data:      parseAPIErrorData(respBody),
		}
	}

	return respBody, nil
}

func (c *httpClient) get(ctx context.Context, path string) ([]byte, error) {
	return c.request(ctx, "GET", path, nil)
}

func (c *httpClient) post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.request(ctx, "POST", path, body)
}

func (c *httpClient) setAPIKey(key string) {
	c.mu.Lock()
	c.apiKey = key
	c.mu.Unlock()
}

func (c *httpClient) setBaseURL(url string) {
	c.mu.Lock()
	c.baseURL = url
	c.mu.Unlock()
}

func (c *httpClient) snapshotConfig() (baseURL string, apiKey string) {
	c.mu.RLock()
	baseURL = c.baseURL
	apiKey = c.apiKey
	c.mu.RUnlock()
	return
}

// parseAPIErrorData extracts all fields from an error response body.
func parseAPIErrorData(body []byte) map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil
	}
	return data
}

func parseAPIErrorMessage(status int, body []byte) string {
	var payload struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Error != "" {
		return payload.Error
	}
	if text := strings.TrimSpace(string(body)); text != "" {
		return text
	}
	return http.StatusText(status)
}
