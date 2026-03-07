package kaizen

import (
	"context"
	"encoding/json"
	"fmt"
)

// AkumaClient is the client for Akuma API.
type AkumaClient struct {
	http *httpClient
}

// Query translates natural language to SQL.
func (c *AkumaClient) Query(ctx context.Context, req *AkumaQueryRequest) (*AkumaQueryResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/query", req)
	if err != nil {
		return nil, err
	}

	var resp AkumaQueryResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma query response: %w", err)
	}
	return &resp, nil
}

// Explain explains a SQL query in plain English.
func (c *AkumaClient) Explain(ctx context.Context, sql string) (*AkumaExplainResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/explain", map[string]string{"sql": sql})
	if err != nil {
		return nil, err
	}

	var resp AkumaExplainResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma explain response: %w", err)
	}
	return &resp, nil
}

// SetSchema persists a manual schema source used by Akuma query generation.
func (c *AkumaClient) SetSchema(ctx context.Context, req *AkumaSchemaRequest) (*AkumaSchemaResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/schema", req)
	if err != nil {
		return nil, err
	}

	var resp AkumaSchemaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma schema response: %w", err)
	}
	return &resp, nil
}

// ListSources lists persisted Akuma data sources.
func (c *AkumaClient) ListSources(ctx context.Context) (*AkumaSourcesResponse, error) {
	data, err := c.http.get(ctx, "/v1/akuma/sources")
	if err != nil {
		return nil, err
	}

	var resp AkumaSourcesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma sources response: %w", err)
	}
	return &resp, nil
}

// CreateSource creates a live Akuma data source.
func (c *AkumaClient) CreateSource(ctx context.Context, req *AkumaCreateSourceRequest) (*AkumaSourceMutationResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/sources", req)
	if err != nil {
		return nil, err
	}

	var resp AkumaSourceMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma source mutation response: %w", err)
	}
	return &resp, nil
}

// DeleteSource deletes a persisted Akuma data source.
func (c *AkumaClient) DeleteSource(ctx context.Context, sourceID string) (*AkumaSourceMutationResponse, error) {
	data, err := c.http.request(ctx, "DELETE", fmt.Sprintf("/v1/akuma/sources/%s", sourceID), nil)
	if err != nil {
		return nil, err
	}

	var resp AkumaSourceMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma source delete response: %w", err)
	}
	return &resp, nil
}

// SyncSource triggers an immediate live sync for a persisted Akuma data source.
func (c *AkumaClient) SyncSource(ctx context.Context, sourceID string) (*AkumaSourceMutationResponse, error) {
	data, err := c.http.post(ctx, fmt.Sprintf("/v1/akuma/sources/%s/sync", sourceID), map[string]any{})
	if err != nil {
		return nil, err
	}

	var resp AkumaSourceMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma source sync response: %w", err)
	}
	return &resp, nil
}
