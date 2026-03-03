# DriftGate Go SDK

Canonical V4 envelope SDK for Go.

See canonical envelope docs: [`docs/sdk/response-envelope.md`](/Users/jordandavis/Documents/Code/DriftGateAI/driftgate-v4-sdk-wt/docs/sdk/response-envelope.md).

## Install

```bash
go get github.com/driftgate/sdk-go@v0.1.0
```

## Hello World (2 lines)

```go
session, _ := driftgatesdk.NewClient("https://api.driftgate.ai").Session.Start(driftgatesdk.SessionStartRequest{Agent: "refund-agent"})
_, _ = session.Execute(driftgatesdk.ExecutionRequest{Input: map[string]any{"orderId": "123"}})
```

## Full Example

```go
client := driftgatesdk.NewClient("https://api.driftgate.ai")
client.BearerToken = "token"
session, _ := client.Session.Start(driftgatesdk.SessionStartRequest{
  Agent: "refund-agent",
  Policy: &driftgatesdk.PolicyRef{Ref: "refund", Version: "2026-02"},
  Route: &driftgatesdk.RouteRef{Provider: "openai", Model: "gpt-4.1-mini", Region: "us-east-1"},
})
resp, _ := session.Execute(driftgatesdk.ExecutionRequest{
  Input: map[string]any{"orderId": "123"},
  Risk: &driftgatesdk.RiskMeta{Decision: "review"},
})
_ = resp.Meta.RequestID
```
