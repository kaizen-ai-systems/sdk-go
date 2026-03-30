package kaizen

// KaizenError is the base SDK error type.
type KaizenError struct {
	Message   string
	Status    int
	Code      string
	RequestID string
	Data      map[string]interface{} // Extra fields from the error response body.
}

func (e *KaizenError) Error() string {
	return e.Message
}

// AuthError represents an authentication failure.
type AuthError struct {
	KaizenError
}

// RateLimitError represents a rate limit response.
type RateLimitError struct {
	KaizenError
	RetryAfter int
}
