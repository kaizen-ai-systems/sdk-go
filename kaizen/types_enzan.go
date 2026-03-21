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

// EnzanModelCostRequest is the request for model-level cost analytics.
type EnzanModelCostRequest struct {
	Window TimeWindow `json:"window"`
}

// EnzanLLMPricingUpsertRequest is the request for LLM pricing upserts.
type EnzanLLMPricingUpsertRequest struct {
	Provider                 string  `json:"provider"`
	Model                    string  `json:"model"`
	DisplayName              string  `json:"display_name,omitempty"`
	InputCostPer1KTokensUSD  float64 `json:"input_cost_per_1k_tokens_usd"`
	OutputCostPer1KTokensUSD float64 `json:"output_cost_per_1k_tokens_usd"`
	Currency                 string  `json:"currency,omitempty"`
	Active                   *bool   `json:"active,omitempty"`
}

// EnzanGPUPricingUpsertRequest is the request for GPU pricing upserts.
type EnzanGPUPricingUpsertRequest struct {
	Provider      string  `json:"provider"`
	GPUType       string  `json:"gpu_type"`
	DisplayName   string  `json:"display_name,omitempty"`
	HourlyRateUSD float64 `json:"hourly_rate_usd"`
	Currency      string  `json:"currency,omitempty"`
	Active        *bool   `json:"active,omitempty"`
}

// EnzanFilters contains filters for Enzan queries.
type EnzanFilters struct {
	Projects  []string `json:"projects,omitempty"`
	Models    []string `json:"models,omitempty"`
	Teams     []string `json:"teams,omitempty"`
	Providers []string `json:"providers,omitempty"`
	Endpoints []string `json:"endpoints,omitempty"`
}

// EnzanSummaryRow represents a row in the summary response.
type EnzanSummaryRow struct {
	Project    string  `json:"project,omitempty"`
	Model      string  `json:"model,omitempty"`
	Team       string  `json:"team,omitempty"`
	Provider   string  `json:"provider,omitempty"`
	Endpoint   string  `json:"endpoint,omitempty"`
	CostUSD    float64 `json:"cost_usd"`
	GPUHours   float64 `json:"gpu_hours"`
	Requests   int64   `json:"requests"`
	TokensIn   int64   `json:"tokens_in"`
	TokensOut  int64   `json:"tokens_out"`
	AvgUtilPct float64 `json:"avg_util_pct,omitempty"`
}

// EnzanSummaryTotal contains aggregate totals for a summary window.
type EnzanSummaryTotal struct {
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
	Total     EnzanSummaryTotal `json:"total"`
	APICosts  *APICostSummary   `json:"apiCosts,omitempty"`
}

// EnzanModelCategoryBreakdown represents prompt complexity distribution per model.
type EnzanModelCategoryBreakdown struct {
	Category        string  `json:"category"`
	Queries         int64   `json:"queries"`
	PromptTokens    int64   `json:"prompt_tokens"`
	OutputTokens    int64   `json:"output_tokens"`
	CostUSD         float64 `json:"cost_usd"`
	Percentage      float64 `json:"percentage"`
	AvgCostPerQuery float64 `json:"avg_cost_per_query"`
}

// EnzanModelCostRow represents a single model row in the cost breakdown response.
type EnzanModelCostRow struct {
	Model           string                        `json:"model"`
	Queries         int64                         `json:"queries"`
	PromptTokens    int64                         `json:"prompt_tokens"`
	OutputTokens    int64                         `json:"output_tokens"`
	CostUSD         float64                       `json:"cost_usd"`
	Percentage      float64                       `json:"percentage"`
	AvgCostPerQuery float64                       `json:"avg_cost_per_query"`
	Categories      []EnzanModelCategoryBreakdown `json:"categories,omitempty"`
}

// EnzanModelCostTotal contains aggregate totals for model-level spend.
type EnzanModelCostTotal struct {
	Queries      int64   `json:"queries"`
	PromptTokens int64   `json:"prompt_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

// EnzanModelCostResponse is the response from model-level cost analytics.
type EnzanModelCostResponse struct {
	Window    string              `json:"window"`
	StartTime string              `json:"startTime"`
	EndTime   string              `json:"endTime"`
	Rows      []EnzanModelCostRow `json:"rows"`
	Total     EnzanModelCostTotal `json:"total"`
}

// EnzanLLMPricing represents one LLM pricing catalog row.
type EnzanLLMPricing struct {
	Provider                 string  `json:"provider"`
	Model                    string  `json:"model"`
	DisplayName              string  `json:"display_name"`
	InputCostPer1KTokensUSD  float64 `json:"input_cost_per_1k_tokens_usd"`
	OutputCostPer1KTokensUSD float64 `json:"output_cost_per_1k_tokens_usd"`
	Currency                 string  `json:"currency"`
	Active                   bool    `json:"active"`
}

// EnzanGPUPricing represents one GPU pricing catalog row.
type EnzanGPUPricing struct {
	Provider      string  `json:"provider"`
	GPUType       string  `json:"gpu_type"`
	DisplayName   string  `json:"display_name"`
	HourlyRateUSD float64 `json:"hourly_rate_usd"`
	Currency      string  `json:"currency"`
	Active        bool    `json:"active"`
}

// EnzanLLMPricingMutationResponse is the upsert response for model pricing.
type EnzanLLMPricingMutationResponse struct {
	Status  string          `json:"status"`
	Pricing EnzanLLMPricing `json:"pricing"`
}

// EnzanGPUPricingMutationResponse is the upsert response for GPU pricing.
type EnzanGPUPricingMutationResponse struct {
	Status  string          `json:"status"`
	Pricing EnzanGPUPricing `json:"pricing"`
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
	Endpoint   string            `json:"endpoint,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
	HourlyRate float64           `json:"hourlyRate"`
	CreatedAt  string            `json:"createdAt,omitempty"`
	LastSeenAt string            `json:"lastSeenAt,omitempty"`
}

// EnzanAlert represents an alert configuration.
type EnzanAlert struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      AlertType         `json:"type"`
	Threshold float64           `json:"threshold"`
	Window    string            `json:"window"`
	Labels    map[string]string `json:"labels,omitempty"`
	Enabled   bool              `json:"enabled"`
}

// EnzanBurnResponse is the response from Enzan burn.
type EnzanBurnResponse struct {
	BurnRateUSDPerHour float64 `json:"burn_rate_usd_per_hour"`
	Timestamp          string  `json:"timestamp"`
}

// EnzanRecommendationType identifies the optimizer rule.
type EnzanRecommendationType string

const (
	EnzanRecModelDowngrade    EnzanRecommendationType = "model_downgrade"
	EnzanRecDuplicateCaching  EnzanRecommendationType = "duplicate_caching"
	EnzanRecSelfHostBreakeven EnzanRecommendationType = "self_host_breakeven"
	EnzanRecSpendAnomaly      EnzanRecommendationType = "spend_anomaly"
	EnzanRecPriceArbitrage    EnzanRecommendationType = "price_arbitrage"
)

// EnzanOptimizeRequest is the request for the optimizer.
type EnzanOptimizeRequest struct {
	Window TimeWindow `json:"window"`
}

// EnzanRecommendation represents a single optimization suggestion.
type EnzanRecommendation struct {
	Type             EnzanRecommendationType `json:"type"`
	Title            string                  `json:"title"`
	Description      string                  `json:"description"`
	EstimatedSavings float64                 `json:"estimatedSavings"`
	Confidence       float64                 `json:"confidence"`
	Suggestion       string                  `json:"suggestion"`
}

// EnzanOptimizeResponse is the response from the optimizer.
// PotentialSavings is a heuristic upper bound; individual recommendations
// may address overlapping spend, so actual realizable savings may be lower.
type EnzanOptimizeResponse struct {
	Window           string                `json:"window"`
	StartTime        string                `json:"startTime"`
	EndTime          string                `json:"endTime"`
	EfficiencyScore  int                   `json:"efficiencyScore"`
	MonthlySpend     float64               `json:"monthlySpend"`
	PotentialSavings float64               `json:"potentialSavings"`
	Recommendations  []EnzanRecommendation `json:"recommendations"`
}
