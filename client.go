package driftgatesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL     string
	APIKey      string
	BearerToken string
	HTTPClient  *http.Client

	Session *SessionService
}

type SessionService struct {
	client *Client
}

type SessionHandle struct {
	client        *Client
	Session       SessionResource
	StartEnvelope Response[SessionStartData]
}

func NewClient(baseURL string) *Client {
	client := &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
	}
	client.Session = &SessionService{client: client}
	return client
}

func (c *Client) Execute(request EphemeralExecuteRequest) (Response[EphemeralExecuteData], error) {
	var envelope Response[EphemeralExecuteData]
	err := c.request(http.MethodPost, "/v4/execute", request, &envelope)
	return envelope, err
}

func (s *SessionService) Start(request SessionStartRequest) (*SessionHandle, error) {
	var envelope Response[SessionStartData]
	if err := s.client.request(http.MethodPost, "/v4/sessions.start", request, &envelope); err != nil {
		return nil, err
	}
	if envelope.Data == nil {
		return nil, fmt.Errorf("v4 session.start returned empty data")
	}
	return &SessionHandle{
		client:        s.client,
		Session:       envelope.Data.Session,
		StartEnvelope: envelope,
	}, nil
}

func (h *SessionHandle) Execute(request ExecutionRequest) (Response[ExecutionResult], error) {
	var envelope Response[ExecutionResult]
	path := fmt.Sprintf("/v4/sessions/%s/executions.execute", h.Session.SessionID)
	err := h.client.request(http.MethodPost, path, request, &envelope)
	return envelope, err
}

func (c *Client) request(method string, path string, payload any, output any) error {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("x-driftgate-api-key", c.APIKey)
	}
	if c.BearerToken != "" {
		req.Header.Set("authorization", "Bearer "+c.BearerToken)
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var envelope Response[map[string]any]
		if err := json.Unmarshal(raw, &envelope); err == nil && envelope.Error != nil {
			envelope.Raw = raw
			return &SDKError{
				Code:      string(envelope.Error.Code),
				Message:   envelope.Error.Message,
				Status:    envelope.Error.Status,
				Retryable: envelope.Error.Retryable,
				RequestID: envelope.Meta.RequestID,
				Details:   envelope.Error.Details,
			}
		}
		return &SDKError{
			Code:      "HTTP_ERROR",
			Message:   fmt.Sprintf("request failed (%d)", resp.StatusCode),
			Status:    resp.StatusCode,
			Retryable: resp.StatusCode >= 500,
		}
	}

	if err := json.Unmarshal(raw, output); err != nil {
		return err
	}

	if typed, ok := output.(*Response[SessionStartData]); ok {
		typed.Raw = raw
	}
	if typed, ok := output.(*Response[ExecutionResult]); ok {
		typed.Raw = raw
	}
	if typed, ok := output.(*Response[EphemeralExecuteData]); ok {
		typed.Raw = raw
	}

	return nil
}
