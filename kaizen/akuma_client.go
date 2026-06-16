package kaizen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AkumaClient is the client for Akuma API.
type AkumaClient struct {
	http *httpClient
}

// Query translates natural language to SQL.
func (c *AkumaClient) Query(ctx context.Context, req *AkumaQueryRequest) (*AkumaQueryResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/query", req)
	if err != nil {
		return nil, err
	}

	var resp AkumaQueryResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma query response: %w", err)
	}
	return &resp, nil
}

// QueryInteractive runs a fresh query through the interactive Akuma protocol.
// To consume a needs_clarification response, use ConsumeClarification.
func (c *AkumaClient) QueryInteractive(ctx context.Context, req *AkumaQueryRequest) (*AkumaInteractiveQueryResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/queries/interactive", req)
	if err != nil {
		return nil, err
	}
	return parseAkumaInteractiveResponse(data)
}

// ConsumeClarification consumes a previously-issued clarification by selecting
// one of the offered options. The Idempotency field is required and is sent as
// the Idempotency-Key header. The first successful consume wins; same-key
// retries replay the persisted result; different-key retries are rejected with
// a 409 KaizenError.
func (c *AkumaClient) ConsumeClarification(ctx context.Context, req *AkumaInteractiveConsumeRequest) (*AkumaInteractiveQueryResponse, error) {
	if req == nil {
		return nil, &KaizenError{Message: "consume request is required", Code: "INVALID_REQUEST"}
	}
	if strings.TrimSpace(req.ClarificationToken) == "" {
		return nil, &KaizenError{Message: "clarificationToken is required", Code: "INVALID_REQUEST"}
	}
	if strings.TrimSpace(req.OptionID) == "" {
		return nil, &KaizenError{Message: "optionId is required", Code: "INVALID_REQUEST"}
	}
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return nil, &KaizenError{Message: "IdempotencyKey is required for consume", Code: "INVALID_REQUEST"}
	}
	body := map[string]string{
		"clarificationToken": req.ClarificationToken,
		"optionId":           req.OptionID,
	}
	headers := map[string]string{"Idempotency-Key": req.IdempotencyKey}
	data, err := c.http.requestWithHeaders(ctx, "POST", "/v1/akuma/queries/interactive", body, headers)
	if err != nil {
		return nil, err
	}
	return parseAkumaInteractiveResponse(data)
}

// parseAkumaInteractiveResponse decodes and validates an interactive envelope.
// Strict shape rules: status required, result must be an object when present,
// completed/rejected require result, needs_clarification requires clarification
// with at least 2 options.
func parseAkumaInteractiveResponse(data []byte) (*AkumaInteractiveQueryResponse, error) {
	var rawEnvelope map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawEnvelope); err != nil || rawEnvelope == nil {
		var topLevel interface{}
		if topErr := json.Unmarshal(data, &topLevel); topErr != nil {
			return nil, &KaizenError{
				Message: fmt.Sprintf("decode interactive akuma query envelope: %v", err),
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		return nil, &KaizenError{
			Message: "interactive query response must be an object",
			Code:    "INVALID_RESPONSE",
			Data:    map[string]interface{}{"response": topLevel},
		}
	}
	var status AkumaInteractiveQueryStatus
	if statusRaw, ok := rawEnvelope["status"]; ok {
		if err := json.Unmarshal(statusRaw, &status); err != nil {
			return nil, &KaizenError{
				Message: fmt.Sprintf("decode interactive akuma query envelope: %v", err),
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
	}
	if strings.TrimSpace(string(status)) == "" {
		return nil, &KaizenError{
			Message: "interactive query response missing status",
			Code:    "INVALID_RESPONSE",
			Data:    parseAPIErrorData(data),
		}
	}
	resultRaw, hasResult := rawEnvelope["result"]
	trimmedResult := bytes.TrimSpace(resultRaw)
	resp := AkumaInteractiveQueryResponse{Status: status, RawResponse: rawEnvelope}
	if len(trimmedResult) > 0 {
		if trimmedResult[0] != '{' {
			return nil, &KaizenError{
				Message: "interactive query response result must be an object",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		var queryResp AkumaQueryResponse
		if err := json.Unmarshal(trimmedResult, &queryResp); err != nil {
			return nil, &KaizenError{
				Message: fmt.Sprintf("decode interactive akuma query result: %v", err),
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		resp.Result = &queryResp
	}
	clarificationRaw, hasClarification := rawEnvelope["clarification"]
	trimmedClarification := bytes.TrimSpace(clarificationRaw)
	// Only strictly parse the clarification payload for the needs_clarification
	// status. A forward-compat future status may ship a permissively-shaped
	// clarification field whose types differ from AkumaClarification; strictly
	// unmarshalling it would break forward compatibility. The raw envelope is
	// always on RawResponse for callers that want to inspect it. This matches
	// the Python SDK's behavior (human review #3 L1).
	if status == AkumaInteractiveQueryStatusNeedsClarification && len(trimmedClarification) > 0 {
		if trimmedClarification[0] != '{' {
			return nil, &KaizenError{
				Message: "interactive query response clarification must be an object",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		var clarification AkumaClarification
		if err := json.Unmarshal(trimmedClarification, &clarification); err != nil {
			return nil, &KaizenError{
				Message: fmt.Sprintf("decode interactive akuma clarification: %v", err),
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		resp.Clarification = &clarification
	}
	// Provenance is only defined for completed results (PR 1c); like
	// clarification above, it is strictly parsed only for the status that
	// defines it so a future status can reshape the field without breaking
	// older SDKs. The raw envelope always carries it for inspection.
	provenanceRaw := bytes.TrimSpace(rawEnvelope["provenance"])
	if status == AkumaInteractiveQueryStatusCompleted && len(provenanceRaw) > 0 && string(provenanceRaw) != "null" {
		if provenanceRaw[0] != '{' {
			return nil, &KaizenError{
				Message: "interactive query response provenance must be an object",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		var provenance AkumaProvenance
		if err := json.Unmarshal(provenanceRaw, &provenance); err != nil {
			return nil, &KaizenError{
				Message: fmt.Sprintf("decode interactive akuma provenance: %v", err),
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		// provenanceFidelity is the load-bearing field — without it a client
		// cannot tell full from limited provenance, so an empty value on a
		// completed envelope is a contract violation, not a tolerable gap.
		if strings.TrimSpace(string(provenance.ProvenanceFidelity)) == "" {
			return nil, &KaizenError{
				Message: "interactive query response provenance missing provenanceFidelity",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		// Contract: arrays are always present (possibly empty), never nil —
		// normalize so the typed payload honors that even if a field was
		// omitted upstream (Claude round-1 M1).
		normalizeAkumaProvenance(&provenance)
		resp.Provenance = &provenance
	}
	if (resp.Status == AkumaInteractiveQueryStatusCompleted || resp.Status == AkumaInteractiveQueryStatusRejected) && !hasResult {
		return nil, &KaizenError{
			Message: "interactive query response missing result",
			Code:    "INVALID_RESPONSE",
			Data:    parseAPIErrorData(data),
		}
	}
	if resp.Status == AkumaInteractiveQueryStatusRejected && (resp.Result == nil || strings.TrimSpace(resp.Result.Error) == "") {
		return nil, &KaizenError{
			Message: "interactive query rejected response missing error",
			Code:    "INVALID_RESPONSE",
			Data:    parseAPIErrorData(data),
		}
	}
	if resp.Status == AkumaInteractiveQueryStatusCompleted && resp.Result != nil && strings.TrimSpace(resp.Result.Error) != "" {
		return nil, &KaizenError{
			Message: "interactive query completed response must not include error",
			Code:    "INVALID_RESPONSE",
			Data:    parseAPIErrorData(data),
		}
	}
	if resp.Status == AkumaInteractiveQueryStatusNeedsClarification {
		if !hasClarification || resp.Clarification == nil {
			return nil, &KaizenError{
				Message: "interactive query needs_clarification response missing clarification",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		if strings.TrimSpace(resp.Clarification.ClarificationToken) == "" {
			return nil, &KaizenError{
				Message: "interactive query needs_clarification response missing clarificationToken",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		if strings.TrimSpace(resp.Clarification.Question) == "" {
			return nil, &KaizenError{
				Message: "interactive query needs_clarification response missing question",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		if len(resp.Clarification.Options) < 2 || len(resp.Clarification.Options) > 4 {
			return nil, &KaizenError{
				Message: "interactive query needs_clarification response requires 2-4 options",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
		for _, opt := range resp.Clarification.Options {
			if strings.TrimSpace(opt.ID) == "" || strings.TrimSpace(opt.Label) == "" {
				return nil, &KaizenError{
					Message: "interactive query needs_clarification option missing id or label",
					Code:    "INVALID_RESPONSE",
					Data:    parseAPIErrorData(data),
				}
			}
		}
		if strings.TrimSpace(resp.Clarification.ExpiresAt) == "" {
			return nil, &KaizenError{
				Message: "interactive query needs_clarification response missing expiresAt",
				Code:    "INVALID_RESPONSE",
				Data:    parseAPIErrorData(data),
			}
		}
	}
	return &resp, nil
}

// Explain explains a SQL query in plain English.
func (c *AkumaClient) Explain(ctx context.Context, sql string) (*AkumaExplainResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/explain", map[string]string{"sql": sql})
	if err != nil {
		return nil, err
	}

	var resp AkumaExplainResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma explain response: %w", err)
	}
	return &resp, nil
}

// SetSchema persists a manual schema source used by Akuma query generation.
func (c *AkumaClient) SetSchema(ctx context.Context, req *AkumaSchemaRequest) (*AkumaSchemaResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/schema", req)
	if err != nil {
		return nil, err
	}

	var resp AkumaSchemaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma schema response: %w", err)
	}
	return &resp, nil
}

// ListSources lists persisted Akuma data sources.
func (c *AkumaClient) ListSources(ctx context.Context) (*AkumaSourcesResponse, error) {
	data, err := c.http.get(ctx, "/v1/akuma/sources")
	if err != nil {
		return nil, err
	}

	var resp AkumaSourcesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma sources response: %w", err)
	}
	return &resp, nil
}

// CreateSource creates a live Akuma data source.
func (c *AkumaClient) CreateSource(ctx context.Context, req *AkumaCreateSourceRequest) (*AkumaSourceMutationResponse, error) {
	data, err := c.http.post(ctx, "/v1/akuma/sources", req)
	if err != nil {
		return nil, err
	}

	var resp AkumaSourceMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma source mutation response: %w", err)
	}
	return &resp, nil
}

// DeleteSource deletes a persisted Akuma data source.
func (c *AkumaClient) DeleteSource(ctx context.Context, sourceID string) (*AkumaSourceMutationResponse, error) {
	data, err := c.http.request(ctx, "DELETE", fmt.Sprintf("/v1/akuma/sources/%s", sourceID), nil)
	if err != nil {
		return nil, err
	}

	var resp AkumaSourceMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma source delete response: %w", err)
	}
	return &resp, nil
}

// SyncSource triggers an immediate live sync for a persisted Akuma data source.
func (c *AkumaClient) SyncSource(ctx context.Context, sourceID string) (*AkumaSourceMutationResponse, error) {
	data, err := c.http.post(ctx, fmt.Sprintf("/v1/akuma/sources/%s/sync", sourceID), map[string]any{})
	if err != nil {
		return nil, err
	}

	var resp AkumaSourceMutationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode akuma source sync response: %w", err)
	}
	return &resp, nil
}
