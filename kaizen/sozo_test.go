package kaizen

import "testing"

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
