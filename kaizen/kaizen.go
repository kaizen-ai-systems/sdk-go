// Package kaizen provides the official Go SDK for Kaizen AI Systems.
// Products: Akuma (NL→SQL) | Enzan (GPU Cost) | Sōzō (Synthetic Data)
package kaizen

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Version is the SDK version
const Version = "1.0.0"

// =============================================================================
// TYPES - AKUMA
// =============================================================================

// SQLDialect represents a database dialect
type SQLDialect string

const (
	DialectPostgres   SQLDialect = "postgres"
	DialectMySQL      SQLDialect = "mysql"
	DialectSnowflake  SQLDialect = "snowflake"
	DialectBigQuery   SQLDialect = "bigquery"
	DialectSQLite     SQLDialect = "sqlite"
	DialectRedshift   SQLDialect = "redshift"
	DialectClickHouse SQLDialect = "clickhouse"
)

// QueryMode represents the query mode for Akuma
type QueryMode string

const (
	ModeSQLOnly       QueryMode = "sql-only"
	ModeSQLAndResults QueryMode = "sql-and-results"
	ModeExplain       QueryMode = "explain"
)

// Guardrails contains security constraints for queries
type Guardrails struct {
	ReadOnly    bool     `json:"readOnly,omitempty"`
	AllowTables []string `json:"allowTables,omitempty"`
	DenyTables  []string `json:"denyTables,omitempty"`
	DenyColumns []string `json:"denyColumns,omitempty"`
	MaxRows     int      `json:"maxRows,omitempty"`
	TimeoutSecs int      `json:"timeoutSecs,omitempty"`
}

// AkumaQueryRequest is the request for Akuma query
type AkumaQueryRequest struct {
	Dialect    SQLDialect  `json:"dialect"`
	Prompt     string      `json:"prompt"`
	Mode       QueryMode   `json:"mode,omitempty"`
	MaxRows    int         `json:"maxRows,omitempty"`
	Guardrails *Guardrails `json:"guardrails,omitempty"`
}

// AkumaQueryResponse is the response from Akuma query
type AkumaQueryResponse struct {
	SQL         string                   `json:"sql"`
	Rows        []map[string]interface{} `json:"rows,omitempty"`
	Explanation string                   `json:"explanation,omitempty"`
	Tables      []string                 `json:"tables,omitempty"`
	Warnings    []string                 `json:"warnings,omitempty"`
	Error       string                   `json:"error,omitempty"`
}

// AkumaExplainResponse is the response from Akuma explain
type AkumaExplainResponse struct {
	SQL         string `json:"sql"`
	Explanation string `json:"explanation"`
}

// =============================================================================
// TYPES - ENZAN
// =============================================================================

// TimeWindow represents a time window for queries
type TimeWindow string

const (
	Window1Hour  TimeWindow = "1h"
	Window24Hour TimeWindow = "24h"
	Window7Day   TimeWindow = "7d"
	Window30Day  TimeWindow = "30d"
)

// GroupByDimension represents a dimension to group by
type GroupByDimension string

const (
	GroupByProject  GroupByDimension = "project"
	GroupByModel    GroupByDimension = "model"
	GroupByTeam     GroupByDimension = "team"
	GroupByProvider GroupByDimension = "provider"
	GroupByEndpoint GroupByDimension = "endpoint"
)

// AlertType represents the type of alert
type AlertType string

const (
	AlertCostThreshold  AlertType = "cost_threshold"
	AlertUsageSpike     AlertType = "usage_spike"
	AlertIdleResource   AlertType = "idle_resource"
	AlertBudgetExceeded AlertType = "budget_exceeded"
)

// EnzanSummaryRequest is the request for Enzan summary
type EnzanSummaryRequest struct {
	Window  TimeWindow         `json:"window"`
	GroupBy []GroupByDimension `json:"groupBy,omitempty"`
	Filters *EnzanFilters      `json:"filters,omitempty"`
}

// EnzanFilters contains filters for Enzan queries
type EnzanFilters struct {
	Projects  []string `json:"projects,omitempty"`
	Models    []string `json:"models,omitempty"`
	Teams     []string `json:"teams,omitempty"`
	Providers []string `json:"providers,omitempty"`
}

// EnzanSummaryRow represents a row in the summary response
type EnzanSummaryRow struct {
	Project   string  `json:"project,omitempty"`
	Model     string  `json:"model,omitempty"`
	Team      string  `json:"team,omitempty"`
	Provider  string  `json:"provider,omitempty"`
	CostUSD   float64 `json:"cost_usd"`
	GPUHours  float64 `json:"gpu_hours"`
	Requests  int64   `json:"requests"`
	TokensIn  int64   `json:"tokens_in"`
	TokensOut int64   `json:"tokens_out"`
}

// EnzanSummaryResponse is the response from Enzan summary
type EnzanSummaryResponse struct {
	Window    string            `json:"window"`
	StartTime string            `json:"startTime"`
	EndTime   string            `json:"endTime"`
	Rows      []EnzanSummaryRow `json:"rows"`
	Total     struct {
		CostUSD  float64 `json:"cost_usd"`
		GPUHours float64 `json:"gpu_hours"`
		Requests int64   `json:"requests"`
	} `json:"total"`
}

// EnzanResource represents a GPU resource
type EnzanResource struct {
	ID         string            `json:"id"`
	Provider   string            `json:"provider"`
	GPUType    string            `json:"gpuType"`
	GPUCount   int               `json:"gpuCount"`
	Region     string            `json:"region,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	HourlyRate float64           `json:"hourlyRate"`
}

// EnzanAlert represents an alert configuration
type EnzanAlert struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      AlertType `json:"type"`
	Threshold float64   `json:"threshold"`
	Window    string    `json:"window"`
	Enabled   bool      `json:"enabled"`
}

// EnzanBurnResponse is the response from Enzan burn
type EnzanBurnResponse struct {
	BurnRateUSDPerHour float64 `json:"burn_rate_usd_per_hour"`
	Timestamp          string  `json:"timestamp"`
}

// =============================================================================
// TYPES - SŌZŌ
// =============================================================================

// CorrelationType represents a correlation type
type CorrelationType string

const (
	CorrelationPositive CorrelationType = "positive"
	CorrelationNegative CorrelationType = "negative"
)

// SozoGenerateRequest is the request for Sozo generate
type SozoGenerateRequest struct {
	Schema       map[string]string          `json:"schema,omitempty"`
	SchemaName   string                     `json:"schemaName,omitempty"`
	Records      int                        `json:"records"`
	Correlations map[string]CorrelationType `json:"correlations,omitempty"`
	Seed         *int                       `json:"seed,omitempty"`
}

// SozoColumnStats represents statistics for a column
type SozoColumnStats struct {
	Type        string         `json:"type"`
	Min         *float64       `json:"min,omitempty"`
	Max         *float64       `json:"max,omitempty"`
	Mean        *float64       `json:"mean,omitempty"`
	UniqueCount *int           `json:"uniqueCount,omitempty"`
	Values      map[string]int `json:"values,omitempty"`
}

// SozoGenerateResponse is the response from Sozo generate
type SozoGenerateResponse struct {
	Columns []string                   `json:"columns"`
	Rows    []map[string]interface{}   `json:"rows"`
	Stats   map[string]SozoColumnStats `json:"stats"`
}

// ToCSV converts the response to CSV format
func (r *SozoGenerateResponse) ToCSV() (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Write header
	if err := w.Write(r.Columns); err != nil {
		return "", err
	}

	// Write rows
	for _, row := range r.Rows {
		record := make([]string, len(r.Columns))
		for i, col := range r.Columns {
			record[i] = fmt.Sprintf("%v", row[col])
		}
		if err := w.Write(record); err != nil {
			return "", err
		}
	}

	w.Flush()
	return buf.String(), w.Error()
}

// ToJSONL converts the response to JSON Lines format
func (r *SozoGenerateResponse) ToJSONL() (string, error) {
	var buf bytes.Buffer
	for i, row := range r.Rows {
		data, err := json.Marshal(row)
		if err != nil {
			return "", err
		}
		buf.Write(data)
		if i < len(r.Rows)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String(), nil
}

// SozoSchemaInfo represents a predefined schema
type SozoSchemaInfo struct {
	Name    string            `json:"name"`
	Columns map[string]string `json:"columns"`
}

// =============================================================================
// ERRORS
// =============================================================================

// KaizenError is the base error type
type KaizenError struct {
	Message string
	Status  int
	Code    string
}

func (e *KaizenError) Error() string {
	return e.Message
}

// AuthError represents an authentication error
type AuthError struct {
	KaizenError
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	KaizenError
	RetryAfter int
}

// =============================================================================
// HTTP CLIENT
// =============================================================================

type httpClient struct {
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
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return nil, &AuthError{KaizenError{Message: "Invalid or missing API key", Status: 401, Code: "AUTH_ERROR"}}
	}
	if resp.StatusCode == 429 {
		return nil, &RateLimitError{KaizenError: KaizenError{Message: "Rate limit exceeded", Status: 429, Code: "RATE_LIMIT"}}
	}
	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		json.Unmarshal(respBody, &errResp)
		return nil, &KaizenError{Message: errResp.Error, Status: resp.StatusCode}
	}

	return respBody, nil
}

func (c *httpClient) get(ctx context.Context, path string) ([]byte, error) {
	return c.request(ctx, "GET", path, nil)
}

func (c *httpClient) post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.request(ctx, "POST", path, body)
}

// =============================================================================
// AKUMA CLIENT
// =============================================================================

// AkumaClient is the client for Akuma API
type AkumaClient struct {
	http *httpClient
}

// Query translates natural language to SQL
func (c *AkumaClient) Query(ctx context.Context, req *AkumaQueryRequest) (*AkumaQueryResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/query", req)
	if err != nil {
		return nil, err
	}
	var resp AkumaQueryResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Explain explains a SQL query in plain English
func (c *AkumaClient) Explain(ctx context.Context, sql string) (*AkumaExplainResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/explain", map[string]string{"sql": sql})
	if err != nil {
		return nil, err
	}
	var resp AkumaExplainResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// =============================================================================
// ENZAN CLIENT
// =============================================================================

// EnzanClient is the client for Enzan API
type EnzanClient struct {
	http *httpClient
}

// Summary gets GPU cost summary for a time window
func (c *EnzanClient) Summary(ctx context.Context, req *EnzanSummaryRequest) (*EnzanSummaryResponse, error) {
	data, err := c.http.post(ctx, "/v1/enzan/summary", req)
	if err != nil {
		return nil, err
	}
	var resp EnzanSummaryResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Burn gets current burn rate
func (c *EnzanClient) Burn(ctx context.Context) (*EnzanBurnResponse, error) {
	data, err := c.http.get(ctx, "/v1/enzan/burn")
	if err != nil {
		return nil, err
	}
	var resp EnzanBurnResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListResources lists registered GPU resources
func (c *EnzanClient) ListResources(ctx context.Context) ([]EnzanResource, error) {
	data, err := c.http.get(ctx, "/v1/enzan/resources")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Resources []EnzanResource `json:"resources"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Resources, nil
}

// RegisterResource registers a GPU resource
func (c *EnzanClient) RegisterResource(ctx context.Context, resource *EnzanResource) error {
	_, err := c.http.post(ctx, "/v1/enzan/resources", resource)
	return err
}

// ListAlerts lists configured alerts
func (c *EnzanClient) ListAlerts(ctx context.Context) ([]EnzanAlert, error) {
	data, err := c.http.get(ctx, "/v1/enzan/alerts")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Alerts []EnzanAlert `json:"alerts"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Alerts, nil
}

// CreateAlert creates an alert
func (c *EnzanClient) CreateAlert(ctx context.Context, alert *EnzanAlert) error {
	_, err := c.http.post(ctx, "/v1/enzan/alerts", alert)
	return err
}

// =============================================================================
// SŌZŌ CLIENT
// =============================================================================

// SozoClient is the client for Sōzō API
type SozoClient struct {
	http *httpClient
}

// Generate generates synthetic data
func (c *SozoClient) Generate(ctx context.Context, req *SozoGenerateRequest) (*SozoGenerateResponse, error) {
	data, err := c.http.post(ctx, "/v1/sozo/generate", req)
	if err != nil {
		return nil, err
	}
	var resp SozoGenerateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListSchemas lists available predefined schemas
func (c *SozoClient) ListSchemas(ctx context.Context) ([]SozoSchemaInfo, error) {
	data, err := c.http.get(ctx, "/v1/sozo/schemas")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Schemas []SozoSchemaInfo `json:"schemas"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Schemas, nil
}

// =============================================================================
// MAIN CLIENT
// =============================================================================

// ClientConfig contains configuration for the client
type ClientConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// Client is the main Kaizen client
type Client struct {
	Akuma *AkumaClient
	Enzan *EnzanClient
	Sozo  *SozoClient
	http  *httpClient
}

// NewClient creates a new Kaizen client
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

// SetAPIKey sets the API key
func (c *Client) SetAPIKey(key string) {
	c.http.apiKey = key
}

// Health checks API health
func (c *Client) Health(ctx context.Context) (map[string]interface{}, error) {
	data, err := c.http.get(ctx, "/health")
	if err != nil {
		return nil, err
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// =============================================================================
// DEFAULT CLIENT
// =============================================================================

var defaultClient = NewClient(nil)

// Akuma returns the default Akuma client
func Akuma() *AkumaClient { return defaultClient.Akuma }

// Enzan returns the default Enzan client
func Enzan() *EnzanClient { return defaultClient.Enzan }

// Sozo returns the default Sōzō client
func Sozo() *SozoClient { return defaultClient.Sozo }

// SetAPIKey sets the API key for the default client
func SetAPIKey(key string) { defaultClient.SetAPIKey(key) }
