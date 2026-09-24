package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		baseURL:    "https://api.deepseek.com/chat/completions",
	}
}

func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	reqPayload := DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: prompt},
		},
		Thinking: ThinkingConfig{
			Type: "disabled",
		},
		Stream: false,
	}

	jsonPayload, err := json.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		errorResponse := DeepSeekErrorResponse{}
		if err := json.Unmarshal(respBody, &errorResponse); err == nil && errorResponse.Error.Message != "" {
			return "", fmt.Errorf("deepseek error: %s, status: %d", errorResponse.Error.Message, resp.StatusCode)
		}

		return "", fmt.Errorf("deepseek error: status %d: %s", resp.StatusCode, string(respBody))
	}

	apiResponse := DeepSeekResponse{}
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return "", err
	}

	if len(apiResponse.Choices) > 0 {
		return apiResponse.Choices[0].Message.Content, nil
	}

	return "", errors.New("deepseek: empty choices")
}
