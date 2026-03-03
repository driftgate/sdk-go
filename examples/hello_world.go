package main

import "github.com/driftgate/sdk-go"

func main() {
	client := driftgatesdk.NewClient("https://api.driftgate.ai")
	session, _ := client.Session.Start(driftgatesdk.SessionStartRequest{Agent: "refund-agent"})
	_, _ = session.Execute(driftgatesdk.ExecutionRequest{Input: map[string]any{"orderId": "123"}})
}
