package driftgatesdk

type CanonicalErrorCode string

const (
	ErrorAuthInvalid      CanonicalErrorCode = "AUTH_INVALID"
	ErrorPolicyDenied     CanonicalErrorCode = "POLICY_DENIED"
	ErrorRiskExceeded     CanonicalErrorCode = "RISK_EXCEEDED"
	ErrorRouteUnavailable CanonicalErrorCode = "ROUTE_UNAVAILABLE"
	ErrorToolBlocked      CanonicalErrorCode = "TOOL_BLOCKED"
	ErrorRateLimited      CanonicalErrorCode = "RATE_LIMITED"
	ErrorTimeout          CanonicalErrorCode = "TIMEOUT"
	ErrorInternal         CanonicalErrorCode = "INTERNAL"
)

type PolicyRef struct {
	Ref     string `json:"ref"`
	Version string `json:"version"`
}

type RouteRef struct {
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
	Region   string `json:"region,omitempty"`
}

type RiskMeta struct {
	Score    *float64 `json:"score,omitempty"`
	Decision string   `json:"decision,omitempty"`
}

type TimingMs struct {
	Total  float64  `json:"total"`
	Policy *float64 `json:"policy,omitempty"`
	Route  *float64 `json:"route,omitempty"`
	Tool   *float64 `json:"tool,omitempty"`
}

type Meta struct {
	RequestID   string     `json:"requestId"`
	SessionID   string     `json:"sessionId,omitempty"`
	ExecutionID string     `json:"executionId,omitempty"`
	LineageID   string     `json:"lineageId,omitempty"`
	Policy      *PolicyRef `json:"policy,omitempty"`
	Route       *RouteRef  `json:"route,omitempty"`
	Risk        *RiskMeta  `json:"risk,omitempty"`
	TimingMs    TimingMs   `json:"timingMs"`
}

type CanonicalError struct {
	Code      CanonicalErrorCode `json:"code"`
	Message   string             `json:"message"`
	Status    int                `json:"status"`
	Retryable bool               `json:"retryable"`
	Details   map[string]any     `json:"details,omitempty"`
}

type Response[T any] struct {
	Ok    bool            `json:"ok"`
	Data  *T              `json:"data"`
	Meta  Meta            `json:"meta"`
	Error *CanonicalError `json:"error"`
	Raw   []byte          `json:"-"`
}

type SessionResource struct {
	SessionID         string         `json:"sessionId"`
	WorkspaceID       string         `json:"workspaceId"`
	Agent             string         `json:"agent"`
	Subject           string         `json:"subject,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	Policy            *PolicyRef     `json:"policy,omitempty"`
	Route             *RouteRef      `json:"route,omitempty"`
	Risk              *RiskMeta      `json:"risk,omitempty"`
	WorkflowVersionID string         `json:"workflowVersionId,omitempty"`
	CreatedAt         string         `json:"createdAt"`
	ExpiresAt         string         `json:"expiresAt,omitempty"`
}

type SessionStartData struct {
	Session SessionResource `json:"session"`
}

type ExecutionResult struct {
	Run                 map[string]any   `json:"run"`
	Approval            map[string]any   `json:"approval,omitempty"`
	Blocked             bool             `json:"blocked"`
	PolicyDecisions     []map[string]any `json:"policyDecisions"`
	EntitlementDecision map[string]any   `json:"entitlementDecision"`
	UsageEntry          map[string]any   `json:"usageEntry"`
	BoundaryDecision    map[string]any   `json:"boundaryDecision,omitempty"`
	FirewallDecision    map[string]any   `json:"firewallDecision,omitempty"`
}

type EphemeralExecuteData struct {
	Session   SessionResource `json:"session"`
	Execution ExecutionResult `json:"execution"`
}

type SessionStartRequest struct {
	WorkspaceID       string         `json:"workspaceId,omitempty"`
	Agent             string         `json:"agent"`
	Subject           string         `json:"subject,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	Policy            *PolicyRef     `json:"policy,omitempty"`
	Route             *RouteRef      `json:"route,omitempty"`
	Risk              *RiskMeta      `json:"risk,omitempty"`
	WorkflowVersionID string         `json:"workflowVersionId,omitempty"`
	ExpiresAt         string         `json:"expiresAt,omitempty"`
}

type ExecutionRequest struct {
	Input             map[string]any `json:"input"`
	Policy            *PolicyRef     `json:"policy,omitempty"`
	Route             *RouteRef      `json:"route,omitempty"`
	Risk              *RiskMeta      `json:"risk,omitempty"`
	WorkflowVersionID string         `json:"workflowVersionId,omitempty"`
}

type EphemeralExecuteRequest struct {
	SessionStartRequest
	Input map[string]any `json:"input"`
}
