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

// AkumaSchemaRequest sets schema context used by Akuma query generation.
type AkumaSchemaRequest struct {
	Version string       `json:"version,omitempty"`
	Tables  []AkumaTable `json:"tables"`
}

// AkumaSchemaResponse is the response from Akuma schema update.
type AkumaSchemaResponse struct {
	Status string `json:"status"`
	Tables int    `json:"tables"`
}
