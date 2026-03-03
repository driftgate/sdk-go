package driftgatesdk

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionStartAndExecuteEnvelope(t *testing.T) {
	var requests []map[string]any

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = r.Body.Close()
		var parsed map[string]any
		_ = json.Unmarshal(body, &parsed)
		requests = append(requests, parsed)
		w.Header().Set("content-type", "application/json")

		switch r.URL.Path {
		case "/v4/sessions.start":
			_, _ = w.Write([]byte(`{"ok":true,"data":{"session":{"sessionId":"sess_123","workspaceId":"ws_123","agent":"refund-agent","createdAt":"2026-02-28T00:00:00.000Z"}},"meta":{"requestId":"req_1","sessionId":"sess_123","timingMs":{"total":12.3}},"error":null}`))
		case "/v4/sessions/sess_123/executions.execute":
			_, _ = w.Write([]byte(`{"ok":true,"data":{"run":{"id":"run_123"},"blocked":false,"policyDecisions":[],"entitlementDecision":{},"usageEntry":{}},"meta":{"requestId":"req_2","executionId":"run_123","lineageId":"run_123","timingMs":{"total":6.1}},"error":null}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"ok":false,"data":null,"meta":{"requestId":"req_404","timingMs":{"total":1}},"error":{"code":"ROUTE_UNAVAILABLE","message":"not found","status":404,"retryable":false}}`))
		}
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	session, err := client.Session.Start(SessionStartRequest{Agent: "refund-agent", WorkspaceID: "ws_123"})
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if session.Session.SessionID != "sess_123" {
		t.Fatalf("unexpected session id: %s", session.Session.SessionID)
	}

	execResp, err := session.Execute(ExecutionRequest{Input: map[string]any{"orderId": "123"}})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if execResp.Meta.ExecutionID != "run_123" {
		t.Fatalf("unexpected execution id: %s", execResp.Meta.ExecutionID)
	}

	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
	if _, ok := requests[0]["workspaceId"]; !ok {
		t.Fatalf("workspaceId missing from start payload")
	}
	if _, ok := requests[1]["input"]; !ok {
		t.Fatalf("input missing from execute payload")
	}
}

func TestCanonicalErrorParsing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"ok":false,"data":null,"meta":{"requestId":"req_denied","timingMs":{"total":2}},"error":{"code":"POLICY_DENIED","message":"denied","status":403,"retryable":false}}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	_, err := client.Session.Start(SessionStartRequest{Agent: "refund-agent"})
	if err == nil {
		t.Fatalf("expected error")
	}

	sdkErr, ok := err.(*SDKError)
	if !ok {
		t.Fatalf("expected SDKError got %T", err)
	}
	if sdkErr.Code != "POLICY_DENIED" || sdkErr.Status != 403 {
		t.Fatalf("unexpected error: %+v", sdkErr)
	}
}
