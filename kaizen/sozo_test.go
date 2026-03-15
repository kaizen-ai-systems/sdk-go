package kaizen

import (
	"encoding/json"
	"testing"
)

func TestSozoGenerateResponse_ToCSV(t *testing.T) {
	response := &SozoGenerateResponse{
		Columns: []string{"name", "note"},
		Rows: []map[string]interface{}{
			{"name": "Alice", "note": "hello,world"},
		},
	}

	csv, err := response.ToCSV()
	if err != nil {
		t.Fatalf("ToCSV returned error: %v", err)
	}

	want := "name,note\nAlice,\"hello,world\"\n"
	if csv != want {
		t.Fatalf("unexpected CSV output:\nwant: %q\n got: %q", want, csv)
	}
}

func TestSozoGenerateResponse_ToCSV_NilBecomesEmpty(t *testing.T) {
	response := &SozoGenerateResponse{
		Columns: []string{"name", "note"},
		Rows: []map[string]interface{}{
			{"name": "Alice", "note": nil},
		},
	}

	csv, err := response.ToCSV()
	if err != nil {
		t.Fatalf("ToCSV returned error: %v", err)
	}

	want := "name,note\nAlice,\n"
	if csv != want {
		t.Fatalf("unexpected CSV output for nil value:\nwant: %q\n got: %q", want, csv)
	}
}

func TestSozoGenerateResponse_ToJSONL(t *testing.T) {
	response := &SozoGenerateResponse{
		Rows: []map[string]interface{}{{"id": 1}, {"id": 2}},
	}

	jsonl, err := response.ToJSONL()
	if err != nil {
		t.Fatalf("ToJSONL returned error: %v", err)
	}

	want := "{\"id\":1}\n{\"id\":2}"
	if jsonl != want {
		t.Fatalf("unexpected JSONL output: want %q, got %q", want, jsonl)
	}
}

func TestSozoGenerateResponse_UnmarshalStatsAndSchemaDescription(t *testing.T) {
	input := []byte(`{
		"columns": ["score"],
		"rows": [{"score": 1}],
		"stats": {
			"score": {
				"type": "float",
				"nullCount": 2,
				"mean": 1.5,
				"stdDev": 0.5
			}
		}
	}`)

	var got SozoGenerateResponse
	if err := json.Unmarshal(input, &got); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if got.Stats["score"].NullCount != 2 {
		t.Fatalf("unexpected null count: %+v", got.Stats["score"])
	}
	if got.Stats["score"].StdDev == nil || *got.Stats["score"].StdDev != 0.5 {
		t.Fatalf("unexpected std dev: %+v", got.Stats["score"])
	}

	schemasInput := []byte(`{"schemas":[{"name":"saas_customers_v1","description":"Preset","columns":{"id":"uuid4"}}]}`)
	var schemas struct {
		Schemas []SozoSchemaInfo `json:"schemas"`
	}
	if err := json.Unmarshal(schemasInput, &schemas); err != nil {
		t.Fatalf("unexpected schema info unmarshal error: %v", err)
	}
	if len(schemas.Schemas) != 1 || schemas.Schemas[0].Description != "Preset" {
		t.Fatalf("unexpected schema info: %+v", schemas.Schemas)
	}
}
