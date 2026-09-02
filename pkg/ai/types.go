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
	"time"
)

// LLMConfig stores resolved configuration for connecting to an LLM provider
type LLMConfig struct {
	API                 string
	Key                 string
	KeyFile             string
	Model               string
	Diagnose            bool
	AutoHeal            bool
	MaxRetries          int
	SaveFixedDockerfile string
	Lint                bool
	Timeout             time.Duration
	OutputFormat        string
	RedactPatterns      []string
}

// BuildErrorContext encapsulates the full context of a failed build step
type BuildErrorContext struct {
	DockerfilePath    string
	DockerfileContent string
	FailedStageIndex  int
	FailedStageName   string
	FailedCommand     string
	ErrorMessage      string
	RecentLogs        string
	BaseOS            string
	TargetArch        string
	TargetOS          string
}

// DiagnosticResult contains parsed root cause and optimization advice
type DiagnosticResult struct {
	RootCause          string   `json:"root_cause"`
	SuggestedFix       string   `json:"suggested_fix"`
	CodeDiff           string   `json:"code_diff"`
	SizeOptimizations  []string `json:"size_optimizations"`
	QualitySuggestions []string `json:"quality_suggestions"`
}

// AutoHealResult contains the patched Dockerfile and unified diff
type AutoHealResult struct {
	PatchedDockerfile string `json:"patched_dockerfile"`
	Explanation       string `json:"explanation"`
	Diff              string `json:"diff"`
}

// ChatMessage represents a single message in an OpenAI-compatible chat payload
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is the standard OpenAI-compatible request payload
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float32       `json:"temperature"`
}

// ChatCompletionChoice is an element in the choices array
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatCompletionResponse is the standard OpenAI-compatible response payload
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Choices []ChatCompletionChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    any    `json:"code"`
	} `json:"error,omitempty"`
}
