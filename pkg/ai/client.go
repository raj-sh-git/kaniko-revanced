/*
Copyright 2026 The kaniko-revanced Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Client is a resilient OpenAI-compatible HTTP client
type Client struct {
	cfg        LLMConfig
	httpClient *http.Client
}

// NewClient initializes an LLM client with resolved configuration
func NewClient(cfg LLMConfig) *Client {
	// Resolve API endpoint
	if cfg.API == "" {
		if val := os.Getenv("KANIKO_LLM_API"); val != "" {
			cfg.API = val
		} else if val := os.Getenv("OPENAI_BASE_URL"); val != "" {
			cfg.API = val
		} else if val := os.Getenv("OPENAI_API_BASE"); val != "" {
			cfg.API = val
		}
	}

	// Resolve API key
	if cfg.Key == "" {
		if cfg.KeyFile != "" {
			if data, err := os.ReadFile(cfg.KeyFile); err == nil {
				cfg.Key = strings.TrimSpace(string(data))
			}
		}
	}
	if cfg.Key == "" {
		if val := os.Getenv("KANIKO_LLM_KEY"); val != "" {
			cfg.Key = val
		} else if val := os.Getenv("OPENAI_API_KEY"); val != "" {
			cfg.Key = val
		} else if val := os.Getenv("GEMINI_API_KEY"); val != "" {
			cfg.Key = val
		} else if val := os.Getenv("DEEPSEEK_API_KEY"); val != "" {
			cfg.Key = val
		}
	}

	// Resolve Model
	if cfg.Model == "" {
		if val := os.Getenv("KANIKO_LLM_MODEL"); val != "" {
			cfg.Model = val
		} else {
			cfg.Model = "gpt-4o-mini"
		}
	}

	// Resolve Timeout
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}

	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// IsConfigured returns true if an API endpoint is provided
func (c *Client) IsConfigured() bool {
	return c.cfg.API != ""
}

// Complete dispatches a chat completion request to the OpenAI-compatible endpoint
func (c *Client) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if !c.IsConfigured() {
		return "", errors.New("no LLM API endpoint configured (use --llm-api or KANIKO_LLM_API)")
	}

	endpoint := c.resolveEndpointURL(c.cfg.API)

	// Sanitize prompts before transmission
	cleanSystem := SanitizeText(systemPrompt, c.cfg.RedactPatterns)
	cleanUser := SanitizeText(userPrompt, c.cfg.RedactPatterns)

	reqPayload := ChatCompletionRequest{
		Model: c.cfg.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: cleanSystem},
			{Role: "user", Content: cleanUser},
		},
		Temperature: 0.2,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", errors.Wrap(err, "failed to serialize LLM request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.cfg.Key != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.Key))
	}

	logrus.Debugf("Sending AI request to %s (model: %s)", endpoint, c.cfg.Model)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", errors.Wrap(err, "LLM request failed")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read LLM response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("LLM API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return "", errors.Wrap(err, "failed to decode LLM response")
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("LLM error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("LLM returned empty choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *Client) resolveEndpointURL(base string) string {
	clean := strings.TrimRight(base, "/")
	if strings.HasSuffix(clean, "/chat/completions") {
		return clean
	}
	return clean + "/chat/completions"
}
