package kaizen

import (
	"context"
	"encoding/json"
	"fmt"
)

// SozoClient is the client for Sōzō API.
type SozoClient struct {
	http *httpClient
}

// Generate generates synthetic data.
func (c *SozoClient) Generate(ctx context.Context, req *SozoGenerateRequest) (*SozoGenerateResponse, error) {
	data, err := c.http.post(ctx, "/v1/sozo/generate", req)
	if err != nil {
		return nil, err
	}

	var resp SozoGenerateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode sozo generate response: %w", err)
	}
	return &resp, nil
}

// ListSchemas lists available predefined schemas.
func (c *SozoClient) ListSchemas(ctx context.Context) ([]SozoSchemaInfo, error) {
	data, err := c.http.get(ctx, "/v1/sozo/schemas")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Schemas []SozoSchemaInfo `json:"schemas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode sozo schemas response: %w", err)
	}
	return resp.Schemas, nil
}
