package kaizen

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

// CorrelationType represents a correlation type.
type CorrelationType string

const (
	CorrelationPositive CorrelationType = "positive"
	CorrelationNegative CorrelationType = "negative"
)

// SozoGenerateRequest is the request for Sozo generate.
type SozoGenerateRequest struct {
	Schema       map[string]string          `json:"schema,omitempty"`
	SchemaName   string                     `json:"schemaName,omitempty"`
	Records      int                        `json:"records"`
	Correlations map[string]CorrelationType `json:"correlations,omitempty"`
	Seed         *int                       `json:"seed,omitempty"`
}

// SozoColumnStats represents statistics for a column.
type SozoColumnStats struct {
	Type        string         `json:"type"`
	Min         *float64       `json:"min,omitempty"`
	Max         *float64       `json:"max,omitempty"`
	Mean        *float64       `json:"mean,omitempty"`
	UniqueCount *int           `json:"uniqueCount,omitempty"`
	Values      map[string]int `json:"values,omitempty"`
}

// SozoGenerateResponse is the response from Sozo generate.
type SozoGenerateResponse struct {
	Columns []string                   `json:"columns"`
	Rows    []map[string]interface{}   `json:"rows"`
	Stats   map[string]SozoColumnStats `json:"stats"`
}

// ToCSV converts the response to CSV format.
func (r *SozoGenerateResponse) ToCSV() (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(r.Columns); err != nil {
		return "", err
	}

	for _, row := range r.Rows {
		record := make([]string, len(r.Columns))
		for i, column := range r.Columns {
			value, ok := row[column]
			if !ok || value == nil {
				record[i] = ""
				continue
			}
			record[i] = fmt.Sprintf("%v", value)
		}
		if err := w.Write(record); err != nil {
			return "", err
		}
	}

	w.Flush()
	return buf.String(), w.Error()
}

// ToJSONL converts the response to JSON Lines format.
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

// SozoSchemaInfo represents a predefined schema.
type SozoSchemaInfo struct {
	Name    string            `json:"name"`
	Columns map[string]string `json:"columns"`
}
