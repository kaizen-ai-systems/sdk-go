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
	AlertCostThreshold         AlertType = "cost_threshold"
	AlertCostAnomaly           AlertType = "cost_anomaly"
	AlertUsageSpike            AlertType = "usage_spike"
	AlertIdleResource          AlertType = "idle_resource"
	AlertBudgetExceeded        AlertType = "budget_exceeded"
	AlertOptimizationAvailable AlertType = "optimization_available"
	AlertPricingChange         AlertType = "pricing_change"
	AlertDailySummary          AlertType = "daily_summary"
)

// CreatableAlertType represents the currently supported alert-create surface.
type CreatableAlertType string

const (
	CreatableAlertCostThreshold         CreatableAlertType = "cost_threshold"
	CreatableAlertCostAnomaly           CreatableAlertType = "cost_anomaly"
	CreatableAlertBudgetExceeded        CreatableAlertType = "budget_exceeded"
	CreatableAlertOptimizationAvailable CreatableAlertType = "optimization_available"
	CreatableAlertPricingChange         CreatableAlertType = "pricing_change"
	CreatableAlertDailySummary          CreatableAlertType = "daily_summary"
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

// EnzanRoutingConfigResponse wraps one user's smart-routing config.
type EnzanRoutingConfigResponse struct {
	Routing EnzanRoutingConfig `json:"routing"`
}

// EnzanRoutingConfig represents one user's smart-routing config.
type EnzanRoutingConfig struct {
	Enabled       bool   `json:"enabled"`
	Provider      string `json:"provider"`
	DefaultModel  string `json:"default_model"`
	SimpleModel   string `json:"simple_model,omitempty"`
	ModerateModel string `json:"moderate_model,omitempty"`
	ComplexModel  string `json:"complex_model,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

// EnzanRoutingConfigUpsertRequest upserts a smart-routing config.
type EnzanRoutingConfigUpsertRequest struct {
	Enabled       *bool   `json:"enabled,omitempty"`
	SimpleModel   *string `json:"simple_model,omitempty"`
	ModerateModel *string `json:"moderate_model,omitempty"`
	ComplexModel  *string `json:"complex_model,omitempty"`
}

// EnzanRoutingConfigMutationResponse is the write response for smart routing.
type EnzanRoutingConfigMutationResponse struct {
	Status  string             `json:"status"`
	Routing EnzanRoutingConfig `json:"routing"`
}

// EnzanRoutingSavingsBreakdown is one smart-routing savings bucket.
type EnzanRoutingSavingsBreakdown struct {
	PromptCategory        string  `json:"prompt_category"`
	OriginalModel         string  `json:"original_model"`
	RoutedModel           string  `json:"routed_model"`
	Queries               int64   `json:"queries"`
	ActualCostUSD         float64 `json:"actual_cost_usd"`
	CounterfactualCostUSD float64 `json:"counterfactual_cost_usd"`
	EstimatedSavingsUSD   float64 `json:"estimated_savings_usd"`
}

// EnzanRoutingSavingsResponse summarizes realized smart-routing savings.
type EnzanRoutingSavingsResponse struct {
	Window                string                         `json:"window"`
	StartTime             string                         `json:"start_time"`
	EndTime               string                         `json:"end_time"`
	Provider              string                         `json:"provider"`
	DefaultModel          string                         `json:"default_model"`
	TotalQueries          int64                          `json:"total_queries"`
	RoutedQueries         int64                          `json:"routed_queries"`
	ActualCostUSD         float64                        `json:"actual_cost_usd"`
	CounterfactualCostUSD float64                        `json:"counterfactual_cost_usd"`
	EstimatedSavingsUSD   float64                        `json:"estimated_savings_usd"`
	Breakdown             []EnzanRoutingSavingsBreakdown `json:"breakdown"`
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

// EnzanPricingRefreshTriggerResponse is the response from POST /v1/enzan/pricing/refresh.
// Status is "queued" (HTTP 202, sweep admitted) or "dropped" (HTTP 429,
// per-process concurrency cap reached; caller should retry shortly).
type EnzanPricingRefreshTriggerResponse struct {
	Status      string `json:"status"`
	TriggeredBy string `json:"triggeredBy"`
}

// EnzanPricingRefreshLogEntry is one row from enzan_pricing_refresh_log.
// Pointer fields are nullable (omitted in JSON when null).
type EnzanPricingRefreshLogEntry struct {
	ID           string  `json:"id"`
	SourceID     *string `json:"sourceId,omitempty"`
	SourceName   *string `json:"sourceName,omitempty"`
	Kind         string  `json:"kind"`
	TriggeredBy  *string `json:"triggeredBy,omitempty"`
	Status       string  `json:"status"`
	RowsUpserted int     `json:"rowsUpserted"`
	RowsSkipped  int     `json:"rowsSkipped"`
	DurationMs   *int32  `json:"durationMs,omitempty"`
	Error        *string `json:"error,omitempty"`
	StartedAt    string  `json:"startedAt"`
	FinishedAt   *string `json:"finishedAt,omitempty"`
}

// EnzanPricingProvider is one row from enzan_pricing_sources (admin view).
type EnzanPricingProvider struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Kind                 string  `json:"kind"`
	Enabled              bool    `json:"enabled"`
	RefreshIntervalHours int     `json:"refreshIntervalHours"`
	HasAdapter           bool    `json:"hasAdapter"`
	LastSuccessAt        *string `json:"lastSuccessAt,omitempty"`
	LastFailureAt        *string `json:"lastFailureAt,omitempty"`
	LastError            *string `json:"lastError,omitempty"`
}

// BoolPtr returns a pointer to b. Convenience constructor for SDK
// payload fields that use *bool to disambiguate "explicit false" from
// "unset" (e.g., EnzanGPUOfferUpsertPayload.TrainingReady).
func BoolPtr(b bool) *bool { return &b }

// Float64Ptr returns a pointer to f. Convenience constructor for SDK
// payload fields that use *float64 to disambiguate "explicit zero"
// (a genuine free offer) from "unset" (a forgotten field that would
// otherwise pass the server's `minimum: 0` validation as a free offer).
func Float64Ptr(f float64) *float64 { return &f }

// EnzanGPUOfferUpsertPayload is the GPU branch of POST /v1/enzan/pricing/offers.
// DeploymentClass enum: on_demand | reserved | spot | committed_monthly.
// InterconnectClass enum: standard | high_speed | infiniband | nvlink | unknown.
//
// TrainingReady and HourlyRateUSD are pointers (not value types) so callers
// can express "unset" distinctly from "explicit false / explicit zero".
// With a plain `float64`, Go's zero value and a deliberate free offer are
// indistinguishable on the wire — a forgetful caller would silently submit
// a $0/hr offer that the server (which accepts `minimum: 0`) would happily
// persist. The pointer forces an explicit decision; the client validates
// non-nil before sending so typos surface as a clear ValueError equivalent
// rather than as a free offer in the catalog.
type EnzanGPUOfferUpsertPayload struct {
	Provider          string   `json:"provider"`
	GPUType           string   `json:"gpuType"`
	DisplayName       string   `json:"displayName"`
	Region            string   `json:"region,omitempty"`
	DeploymentClass   string   `json:"deploymentClass,omitempty"`
	CommitmentTerm    string   `json:"commitmentTerm,omitempty"`
	ClusterSizeMin    int      `json:"clusterSizeMin,omitempty"`
	ClusterSizeMax    int      `json:"clusterSizeMax,omitempty"`
	InterconnectClass string   `json:"interconnectClass,omitempty"`
	TrainingReady     *bool    `json:"trainingReady,omitempty"`
	HourlyRateUSD     *float64 `json:"hourlyRateUSD,omitempty"`
	Currency          string   `json:"currency,omitempty"`
	CurrencyFxNote    string   `json:"currencyFxNote,omitempty"`
	SourceURL         string   `json:"sourceUrl,omitempty"`
}

// EnzanLLMOfferUpsertPayload is the LLM branch of POST /v1/enzan/pricing/offers.
// Cost fields are pointers for the same "unset != explicit zero" reason
// documented on EnzanGPUOfferUpsertPayload.
type EnzanLLMOfferUpsertPayload struct {
	Provider                 string   `json:"provider"`
	Model                    string   `json:"model"`
	DisplayName              string   `json:"displayName"`
	Region                   string   `json:"region,omitempty"`
	CommitmentTerm           string   `json:"commitmentTerm,omitempty"`
	InputCostPer1KTokensUSD  *float64 `json:"inputCostPer1KTokensUSD,omitempty"`
	OutputCostPer1KTokensUSD *float64 `json:"outputCostPer1KTokensUSD,omitempty"`
	Currency                 string   `json:"currency,omitempty"`
	CurrencyFxNote           string   `json:"currencyFxNote,omitempty"`
	SourceURL                string   `json:"sourceUrl,omitempty"`
}

// EnzanPricingOfferUpsertRequest is the body for POST /v1/enzan/pricing/offers.
// Exactly one of GPU or LLM must be set.
type EnzanPricingOfferUpsertRequest struct {
	GPU *EnzanGPUOfferUpsertPayload `json:"gpu,omitempty"`
	LLM *EnzanLLMOfferUpsertPayload `json:"llm,omitempty"`
}

// EnzanGPUOffer is the persisted GPU offer (admin or adapter-sourced).
type EnzanGPUOffer struct {
	ID                string  `json:"id"`
	Provider          string  `json:"provider"`
	GPUType           string  `json:"gpuType"`
	DisplayName       string  `json:"displayName"`
	Region            *string `json:"region,omitempty"`
	DeploymentClass   string  `json:"deploymentClass"`
	CommitmentTerm    *string `json:"commitmentTerm,omitempty"`
	ClusterSizeMin    int     `json:"clusterSizeMin"`
	ClusterSizeMax    *int32  `json:"clusterSizeMax,omitempty"`
	InterconnectClass string  `json:"interconnectClass"`
	TrainingReady     bool    `json:"trainingReady"`
	HourlyRateUSD     float64 `json:"hourlyRateUSD"`
	Currency          string  `json:"currency"`
	CurrencyFxNote    *string `json:"currencyFxNote,omitempty"`
	SourceType        string  `json:"sourceType"`
	SourceID          *string `json:"sourceId,omitempty"`
	SourceURL         *string `json:"sourceUrl,omitempty"`
	SourceFingerprint *string `json:"sourceFingerprint,omitempty"`
	TrustStatus       string  `json:"trustStatus"`
	FetchedAt         string  `json:"fetchedAt"`
	FirstSeenAt       string  `json:"firstSeenAt"`
	LastSeenAt        string  `json:"lastSeenAt"`
	Active            bool    `json:"active"`
}

// EnzanLLMOffer is the persisted LLM offer (admin or adapter-sourced).
type EnzanLLMOffer struct {
	ID                       string  `json:"id"`
	Provider                 string  `json:"provider"`
	Model                    string  `json:"model"`
	DisplayName              string  `json:"displayName"`
	Region                   *string `json:"region,omitempty"`
	CommitmentTerm           *string `json:"commitmentTerm,omitempty"`
	InputCostPer1KTokensUSD  float64 `json:"inputCostPer1KTokensUSD"`
	OutputCostPer1KTokensUSD float64 `json:"outputCostPer1KTokensUSD"`
	Currency                 string  `json:"currency"`
	CurrencyFxNote           *string `json:"currencyFxNote,omitempty"`
	SourceType               string  `json:"sourceType"`
	SourceID                 *string `json:"sourceId,omitempty"`
	SourceURL                *string `json:"sourceUrl,omitempty"`
	SourceFingerprint        *string `json:"sourceFingerprint,omitempty"`
	TrustStatus              string  `json:"trustStatus"`
	FetchedAt                string  `json:"fetchedAt"`
	FirstSeenAt              string  `json:"firstSeenAt"`
	LastSeenAt               string  `json:"lastSeenAt"`
	Active                   bool    `json:"active"`
}

// EnzanPricingOfferUpsertResponse is the response from POST /v1/enzan/pricing/offers.
// Status is "upserted" (HTTP 201) or "stale" (HTTP 409 — newer fetched_at row exists).
type EnzanPricingOfferUpsertResponse struct {
	Status string         `json:"status"`
	GPU    *EnzanGPUOffer `json:"gpu,omitempty"`
	LLM    *EnzanLLMOffer `json:"llm,omitempty"`
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
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Type            AlertType         `json:"type"`
	Threshold       float64           `json:"threshold"`
	Window          string            `json:"window"`
	Labels          map[string]string `json:"labels,omitempty"`
	Enabled         bool              `json:"enabled"`
	EvaluationState string            `json:"evaluationState,omitempty"`
	NextEligibleAt  string            `json:"nextEligibleAt,omitempty"`
	StatusReason    string            `json:"statusReason,omitempty"`
}

// EnzanCreateAlertRequest is the request for creating an alert.
type EnzanCreateAlertRequest struct {
	ID        string             `json:"id,omitempty"`
	Name      string             `json:"name"`
	Type      CreatableAlertType `json:"type"`
	Threshold *float64           `json:"threshold,omitempty"`
	// Window is required when Type is cost_threshold or cost_anomaly, defaults to 30d for
	// optimization_available, must be omitted or set to 24h for daily_summary,
	// and is ignored for pricing_change.
	Window  string            `json:"window,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
	Enabled *bool             `json:"enabled,omitempty"`
}

// EnzanUpdateAlertRequest is the request for updating an alert.
type EnzanUpdateAlertRequest struct {
	Name      *string            `json:"name,omitempty"`
	Threshold *float64           `json:"threshold,omitempty"`
	Window    *string            `json:"window,omitempty"`
	Labels    *map[string]string `json:"labels,omitempty"`
	Enabled   *bool              `json:"enabled,omitempty"`
}

// StatusWithIDResponse is the generic mutation response for created resources.
type StatusWithIDResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

// EnzanAlertMutationResponse is the update response for one alert.
type EnzanAlertMutationResponse struct {
	Status string     `json:"status"`
	Alert  EnzanAlert `json:"alert"`
}

// EnzanAlertEndpoint represents one webhook delivery endpoint for Enzan alerts.
type EnzanAlertEndpoint struct {
	ID               string `json:"id"`
	Kind             string `json:"kind"`
	TargetURL        string `json:"targetUrl"`
	HasSigningSecret bool   `json:"hasSigningSecret"`
	Enabled          bool   `json:"enabled"`
	LastUsedAt       string `json:"lastUsedAt,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// EnzanAlertEndpointCreateRequest is the request for creating a webhook endpoint.
type EnzanAlertEndpointCreateRequest struct {
	TargetURL     string `json:"targetUrl"`
	SigningSecret string `json:"signingSecret,omitempty"`
}

// EnzanAlertEndpointUpdateRequest is the request for updating a webhook endpoint.
type EnzanAlertEndpointUpdateRequest struct {
	TargetURL     *string `json:"targetUrl,omitempty"`
	SigningSecret *string `json:"signingSecret,omitempty"`
	Enabled       *bool   `json:"enabled,omitempty"`
}

// EnzanAlertEndpointMutationResponse is the create response for one endpoint.
type EnzanAlertEndpointMutationResponse struct {
	Status   string             `json:"status"`
	Endpoint EnzanAlertEndpoint `json:"endpoint"`
}

// EnzanAlertEvent represents one fired alert event.
type EnzanAlertEvent struct {
	ID          string                 `json:"id"`
	RuleID      string                 `json:"ruleId,omitempty"`
	Type        AlertType              `json:"type"`
	DedupeKey   string                 `json:"dedupeKey"`
	Payload     map[string]interface{} `json:"payload"`
	TriggeredAt string                 `json:"triggeredAt"`
}

// EnzanAlertDelivery represents one delivery attempt/status row for an alert event.
type EnzanAlertDelivery struct {
	ID               string `json:"id"`
	EventID          string `json:"eventId"`
	EndpointID       string `json:"endpointId,omitempty"`
	Status           string `json:"status"`
	RetryCount       int    `json:"retryCount"`
	NextRetryAt      string `json:"nextRetryAt"`
	LastAttemptedAt  string `json:"lastAttemptedAt,omitempty"`
	LastResponseCode *int   `json:"lastResponseCode,omitempty"`
	LastError        string `json:"lastError,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
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

// EnzanChatRequest is the request for conversational AI cost Q&A.
type EnzanChatRequest struct {
	Message        string `json:"message"`
	ConversationID string `json:"conversationId,omitempty"`
	Window         string `json:"window,omitempty"`
}

// EnzanSuggestedAction is a typed action chip from the chat response.
type EnzanSuggestedAction struct {
	Type   string `json:"type"` // set_window, view_summary, view_costs_by_model, view_optimizer, view_model_pricing, view_gpu_pricing
	Label  string `json:"label"`
	Window string `json:"window,omitempty"`
	Model  string `json:"model,omitempty"`
}

// EnzanChatResponse is the response from the chat endpoint.
type EnzanChatResponse struct {
	ConversationID   string                 `json:"conversationId"`
	Message          string                 `json:"message"`
	EffectiveWindow  string                 `json:"effectiveWindow,omitempty"`
	SuggestedActions []EnzanSuggestedAction `json:"suggestedActions"`
	SupportingData   map[string]any         `json:"supportingData,omitempty"`
}
