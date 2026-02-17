package kaizen

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ClientConfig contains configuration for the client.
type ClientConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// Client is the main Kaizen client.
type Client struct {
	Akuma *AkumaClient
	Enzan *EnzanClient
	Sozo  *SozoClient
	http  *httpClient
}

// NewClient creates a new Kaizen client.
func NewClient(cfg *ClientConfig) *Client {
	if cfg == nil {
		cfg = &ClientConfig{}
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.kaizenaisystems.com"
	}
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("KAIZEN_API_KEY")
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	http := newHTTPClient(cfg.BaseURL, cfg.APIKey, cfg.Timeout)
	return &Client{
		Akuma: &AkumaClient{http: http},
		Enzan: &EnzanClient{http: http},
		Sozo:  &SozoClient{http: http},
		http:  http,
	}
}

// SetAPIKey sets the API key.
func (c *Client) SetAPIKey(key string) {
	c.http.setAPIKey(key)
}

// SetBaseURL updates the API base URL.
func (c *Client) SetBaseURL(url string) {
	c.http.setBaseURL(url)
}

// Health checks API health.
func (c *Client) Health(ctx context.Context) (map[string]interface{}, error) {
	data, err := c.http.get(ctx, "/health")
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode health response: %w", err)
	}
	return resp, nil
}

var defaultClient = NewClient(nil)

// Akuma returns the default Akuma client.
func Akuma() *AkumaClient { return defaultClient.Akuma }

// Enzan returns the default Enzan client.
func Enzan() *EnzanClient { return defaultClient.Enzan }

// Sozo returns the default Sōzō client.
func Sozo() *SozoClient { return defaultClient.Sozo }

// SetAPIKey sets the API key for the default client.
func SetAPIKey(key string) { defaultClient.SetAPIKey(key) }

// SetBaseURL sets the API base URL for the default client.
func SetBaseURL(url string) { defaultClient.SetBaseURL(url) }
