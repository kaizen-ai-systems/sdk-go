package kaizen

import "encoding/json"

// SQLDialect represents a database dialect.
type SQLDialect string

const (
	DialectPostgres  SQLDialect = "postgres"
	DialectMySQL     SQLDialect = "mysql"
	DialectSnowflake SQLDialect = "snowflake"
	DialectBigQuery  SQLDialect = "bigquery"
)

// QueryMode represents the query mode for Akuma.
type QueryMode string

const (
	ModeSQLOnly       QueryMode = "sql-only"
	ModeSQLAndResults QueryMode = "sql-and-results"
	ModeExplain       QueryMode = "explain"
)

// AkumaSourceStatus is the persisted source sync state.
type AkumaSourceStatus string

const (
	AkumaSourceStatusSyncing        AkumaSourceStatus = "syncing"
	AkumaSourceStatusActive         AkumaSourceStatus = "active"
	AkumaSourceStatusError          AkumaSourceStatus = "error"
	AkumaSourceStatusSchemaTooLarge AkumaSourceStatus = "schema_too_large"
)

// Guardrails contains security constraints for queries.
type Guardrails struct {
	ReadOnly    bool     `json:"readOnly,omitempty"`
	AllowTables []string `json:"allowTables,omitempty"`
	DenyTables  []string `json:"denyTables,omitempty"`
	DenyColumns []string `json:"denyColumns,omitempty"`
	MaxRows     int      `json:"maxRows,omitempty"`
	TimeoutSecs int      `json:"timeoutSecs,omitempty"`
}

// AkumaQueryRequest is the request for Akuma query.
type AkumaQueryRequest struct {
	Dialect    SQLDialect  `json:"dialect"`
	Prompt     string      `json:"prompt"`
	Mode       QueryMode   `json:"mode,omitempty"`
	MaxRows    int         `json:"maxRows,omitempty"`
	SourceID   string      `json:"sourceId,omitempty"`
	Guardrails *Guardrails `json:"guardrails,omitempty"`
}

// AkumaQueryResponse is the response from Akuma query.
type AkumaQueryResponse struct {
	SQL         string                   `json:"sql"`
	Rows        []map[string]interface{} `json:"rows,omitempty"`
	Explanation string                   `json:"explanation,omitempty"`
	Tables      []string                 `json:"tables,omitempty"`
	Warnings    []string                 `json:"warnings,omitempty"`
	Error       string                   `json:"error,omitempty"`
}

// AkumaInteractiveQueryStatus is the state returned by the interactive query protocol.
// The status field is x-extensible-enum on the server; clients should treat
// unknown values as forward-compat passthrough using RawResponse.
type AkumaInteractiveQueryStatus string

const (
	AkumaInteractiveQueryStatusCompleted          AkumaInteractiveQueryStatus = "completed"
	AkumaInteractiveQueryStatusRejected           AkumaInteractiveQueryStatus = "rejected"
	AkumaInteractiveQueryStatusNeedsClarification AkumaInteractiveQueryStatus = "needs_clarification"
)

// AkumaClarificationOption is one disambiguation choice presented in a
// needs_clarification response. ID is the stable identifier the caller passes
// back as OptionID on the consume request.
type AkumaClarificationOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// AkumaClarification is the structured payload returned alongside a
// needs_clarification status. ClarificationToken is opaque; clients must echo
// it verbatim on the consume call.
type AkumaClarification struct {
	ClarificationToken string                     `json:"clarificationToken"`
	Question           string                     `json:"question"`
	Options            []AkumaClarificationOption `json:"options"`
	ExpiresAt          string                     `json:"expiresAt"`
}

// AkumaInteractiveQueryResponse is the response from the interactive Akuma
// query protocol. Either Result (completed/rejected) or Clarification
// (needs_clarification) is set; future statuses pass through RawResponse only.
// Provenance is present only on completed results that carry executed SQL.
type AkumaInteractiveQueryResponse struct {
	Status        AkumaInteractiveQueryStatus `json:"status"`
	Result        *AkumaQueryResponse         `json:"result,omitempty"`
	Clarification *AkumaClarification         `json:"clarification,omitempty"`
	Provenance    *AkumaProvenance            `json:"provenance,omitempty"`
	RawResponse   map[string]json.RawMessage  `json:"-"`
}

// AkumaProvenanceFidelity reports how the provenance facts were derived. The
// field is x-extensible-enum on the server; treat unknown values as limited.
type AkumaProvenanceFidelity string

const (
	// AkumaProvenanceFidelityFull: parser-backed extraction succeeded
	// (postgres/mysql) — joins, filters, projections, and aggregations are
	// deterministic facts of the executed SQL.
	AkumaProvenanceFidelityFull AkumaProvenanceFidelity = "full"
	// AkumaProvenanceFidelityLimited: generation-metadata fallback (other
	// dialects, or parser failure). Parsed-fact arrays are empty; tables come
	// from generation metadata. Do not present as equivalent to full.
	AkumaProvenanceFidelityLimited AkumaProvenanceFidelity = "limited"
)

// AkumaProvenance is the structured, server-derived explanation of how a
// completed result was produced. Arrays are always present (possibly empty).
type AkumaProvenance struct {
	ProvenanceFidelity AkumaProvenanceFidelity     `json:"provenanceFidelity"`
	Dialect            string                      `json:"dialect"`
	Tables             []string                    `json:"tables"`
	Joins              []AkumaProvenanceJoin       `json:"joins"`
	Filters            []AkumaProvenanceFilter     `json:"filters"`
	Projections        []AkumaProvenanceProjection `json:"projections"`
	Aggregations       AkumaProvenanceAggregations `json:"aggregations"`
	Warnings           []AkumaProvenanceWarning    `json:"warnings"`
}

// AkumaProvenanceJoin is one join recovered from the executed SQL.
type AkumaProvenanceJoin struct {
	Type      string `json:"type"`
	Left      string `json:"left,omitempty"`
	Right     string `json:"right,omitempty"`
	Condition string `json:"condition,omitempty"`
}

// AkumaProvenanceFilter is one WHERE/HAVING predicate recovered from the
// executed SQL.
type AkumaProvenanceFilter struct {
	Clause     string `json:"clause"`
	Expression string `json:"expression"`
}

// AkumaProvenanceProjection is one select-list entry recovered from the
// executed SQL.
type AkumaProvenanceProjection struct {
	Expression string `json:"expression"`
	Alias      string `json:"alias,omitempty"`
}

// AkumaProvenanceAggregations carries the grouping shape of the executed SQL.
type AkumaProvenanceAggregations struct {
	GroupBy   []string `json:"groupBy"`
	Functions []string `json:"functions"`
}

// AkumaProvenanceWarning is an enumerated server-side provenance warning.
// Known categories: parser_extraction_failed, dialect_provenance_limited,
// statement_shape_unsupported; the set is additive over time.
type AkumaProvenanceWarning struct {
	Category string `json:"category"`
	Detail   string `json:"detail,omitempty"`
}

// normalizeAkumaProvenance enforces the contract's "arrays are always present
// (possibly empty)" shape on a decoded payload, so callers can range/index
// without nil checks.
func normalizeAkumaProvenance(p *AkumaProvenance) {
	if p.Tables == nil {
		p.Tables = []string{}
	}
	if p.Joins == nil {
		p.Joins = []AkumaProvenanceJoin{}
	}
	if p.Filters == nil {
		p.Filters = []AkumaProvenanceFilter{}
	}
	if p.Projections == nil {
		p.Projections = []AkumaProvenanceProjection{}
	}
	if p.Warnings == nil {
		p.Warnings = []AkumaProvenanceWarning{}
	}
	if p.Aggregations.GroupBy == nil {
		p.Aggregations.GroupBy = []string{}
	}
	if p.Aggregations.Functions == nil {
		p.Aggregations.Functions = []string{}
	}
}

// AkumaInteractiveConsumeRequest carries the consume-side fields. IdempotencyKey
// MUST be a non-empty string and is sent as the Idempotency-Key header. The
// field name mirrors the camelCase / snake_case naming in the TypeScript and
// Python SDKs (idempotencyKey / idempotency_key) for cross-SDK symmetry.
type AkumaInteractiveConsumeRequest struct {
	ClarificationToken string `json:"clarificationToken"`
	OptionID           string `json:"optionId"`
	IdempotencyKey     string `json:"-"`
}

// AkumaExplainResponse is the response from Akuma explain.
type AkumaExplainResponse struct {
	SQL         string `json:"sql"`
	Explanation string `json:"explanation"`
}

// AkumaColumn is a schema column definition for Akuma context.
type AkumaColumn struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Nullable    bool     `json:"nullable,omitempty"`
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

// AkumaForeignKey is a schema foreign key definition for Akuma context.
type AkumaForeignKey struct {
	Columns    []string `json:"columns"`
	RefTable   string   `json:"refTable"`
	RefColumns []string `json:"refColumns"`
}

// AkumaTable is a schema table definition for Akuma context.
type AkumaTable struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Columns     []AkumaColumn     `json:"columns"`
	PrimaryKey  []string          `json:"primaryKey,omitempty"`
	ForeignKeys []AkumaForeignKey `json:"foreignKeys,omitempty"`
}

// AkumaSchemaRequest persists a manual schema source.
type AkumaSchemaRequest struct {
	SourceID string       `json:"sourceId,omitempty"`
	Name     string       `json:"name,omitempty"`
	Dialect  SQLDialect   `json:"dialect"`
	Version  string       `json:"version,omitempty"`
	Tables   []AkumaTable `json:"tables"`
}

// AkumaCreateSourceRequest creates a live source.
type AkumaCreateSourceRequest struct {
	Name             string     `json:"name"`
	Dialect          SQLDialect `json:"dialect"`
	TargetSchemas    []string   `json:"targetSchemas,omitempty"`
	ConnectionString string     `json:"connectionString"`
}

// AkumaSource is a persisted Akuma source.
type AkumaSource struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Dialect       SQLDialect        `json:"dialect"`
	IsManual      bool              `json:"isManual"`
	TargetSchemas []string          `json:"targetSchemas"`
	Status        AkumaSourceStatus `json:"status"`
	LastError     string            `json:"lastError,omitempty"`
	LastSyncedAt  string            `json:"lastSyncedAt,omitempty"`
	CreatedAt     string            `json:"createdAt"`
	UpdatedAt     string            `json:"updatedAt"`
}

// AkumaSourcesResponse is the response from list sources.
type AkumaSourcesResponse struct {
	Sources []AkumaSource `json:"sources"`
}

// AkumaSourceMutationResponse is the response from source mutations.
type AkumaSourceMutationResponse struct {
	Status   string       `json:"status"`
	SourceID string       `json:"sourceId,omitempty"`
	Source   *AkumaSource `json:"source,omitempty"`
	Tables   int          `json:"tables,omitempty"`
}

// AkumaSchemaResponse aliases the manual schema mutation response.
type AkumaSchemaResponse = AkumaSourceMutationResponse
