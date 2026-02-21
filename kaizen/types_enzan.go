package kaizen

// TimeWindow represents a time window for queries.
type TimeWindow string

const (
	Window1Hour  TimeWindow = "1h"
	Window24Hour TimeWindow = "24h"
	Window7Day   TimeWindow = "7d"
	Window30Day  TimeWindow = "30d"
)

// GroupByDimension represents a dimension to group by.
type GroupByDimension string

const (
	GroupByProject  GroupByDimension = "project"
	GroupByModel    GroupByDimension = "model"
	GroupByTeam     GroupByDimension = "team"
	GroupByProvider GroupByDimension = "provider"
	GroupByEndpoint GroupByDimension = "endpoint"
)

// AlertType represents the type of alert.
type AlertType string

const (
	AlertCostThreshold  AlertType = "cost_threshold"
	AlertCostAnomaly    AlertType = "cost_anomaly"
	AlertUsageSpike     AlertType = "usage_spike"
	AlertIdleResource   AlertType = "idle_resource"
	AlertBudgetExceeded AlertType = "budget_exceeded"
)

// EnzanSummaryRequest is the request for Enzan summary.
type EnzanSummaryRequest struct {
	Window  TimeWindow         `json:"window"`
	GroupBy []GroupByDimension `json:"groupBy,omitempty"`
	Filters *EnzanFilters      `json:"filters,omitempty"`
}

// EnzanFilters contains filters for Enzan queries.
type EnzanFilters struct {
	Projects  []string `json:"projects,omitempty"`
	Models    []string `json:"models,omitempty"`
	Teams     []string `json:"teams,omitempty"`
	Providers []string `json:"providers,omitempty"`
}

// EnzanSummaryRow represents a row in the summary response.
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

// EnzanSummaryResponse is the response from Enzan summary.
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
	APICosts *APICostSummary `json:"apiCosts,omitempty"`
}

// APICostSummary represents estimated Akuma API token spend.
type APICostSummary struct {
	TotalCostUSD float64 `json:"totalCostUsd"`
	PromptTokens int64   `json:"promptTokens"`
	OutputTokens int64   `json:"outputTokens"`
	Queries      int64   `json:"queries"`
}

// EnzanResource represents a GPU resource.
type EnzanResource struct {
	ID         string            `json:"id"`
	Provider   string            `json:"provider"`
	GPUType    string            `json:"gpuType"`
	GPUCount   int               `json:"gpuCount"`
	Region     string            `json:"region,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	HourlyRate float64           `json:"hourlyRate"`
}

// EnzanAlert represents an alert configuration.
type EnzanAlert struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      AlertType `json:"type"`
	Threshold float64   `json:"threshold"`
	Window    string    `json:"window"`
	Enabled   bool      `json:"enabled"`
}

// EnzanBurnResponse is the response from Enzan burn.
type EnzanBurnResponse struct {
	BurnRateUSDPerHour float64 `json:"burn_rate_usd_per_hour"`
	Timestamp          string  `json:"timestamp"`
}
