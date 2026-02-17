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

// SetSchema updates schema context used by Akuma query generation.
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
