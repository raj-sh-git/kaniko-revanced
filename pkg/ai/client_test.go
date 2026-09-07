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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Complete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp := ChatCompletionResponse{
			Choices: []ChatCompletionChoice{
				{
					Message: ChatMessage{
						Role:    "assistant",
						Content: "Simulated LLM diagnostic response",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(LLMConfig{
		API:     server.URL,
		Key:     "test-key",
		Model:   "test-model",
		Timeout: 5 * time.Second,
	})

	res, err := client.Complete(context.Background(), "System prompt", "User prompt")
	if err != nil {
		t.Fatalf("unexpected Complete error: %v", err)
	}

	if res != "Simulated LLM diagnostic response" {
		t.Errorf("expected 'Simulated LLM diagnostic response', got %q", res)
	}
}

func TestClient_CheckReachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Reachable case
	client := NewClient(LLMConfig{
		API:     server.URL,
		Key:     "test-key",
		Timeout: 5 * time.Second,
	})
	if err := client.CheckReachable(context.Background()); err != nil {
		t.Fatalf("expected reachable, got error: %v", err)
	}

	// Unreachable case (closed server)
	server.Close()
	if err := client.CheckReachable(context.Background()); err == nil {
		t.Fatalf("expected unreachable error, got nil")
	}

	// Unconfigured client case
	unconfiguredClient := NewClient(LLMConfig{})
	if err := unconfiguredClient.CheckReachable(context.Background()); err == nil {
		t.Fatalf("expected error for unconfigured client, got nil")
	}
}
