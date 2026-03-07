package kaizen

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
