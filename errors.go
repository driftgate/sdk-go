package driftgatesdk

type SDKError struct {
	Code      string
	Message   string
	Status    int
	Retryable bool
	RequestID string
	Details   map[string]any
}

func (e *SDKError) Error() string {
	return e.Code + ": " + e.Message
}
