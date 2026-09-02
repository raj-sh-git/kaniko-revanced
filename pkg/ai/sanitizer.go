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
	"regexp"
)

var (
	awsAccessKeyRegex = regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`)
	githubTokenRegex  = regexp.MustCompile(`\b(ghp_[0-9a-zA-Z]{36}|github_pat_[0-9a-zA-Z_]{50,90})\b`)
	jwtRegex          = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`)
	bearerRegex       = regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9\-._~+/]+=*`)
	privateKeyRegex   = regexp.MustCompile(`(?s)-----BEGIN [A-Z ]+PRIVATE KEY-----.*?-----END [A-Z ]+PRIVATE KEY-----`)
	urlPassRegex      = regexp.MustCompile(`(https?://[^:\s]+):([^@\s]+)@`)
	envSecretRegex    = regexp.MustCompile(`(?i)(password|secret|api_key|token|auth|credentials?)\s*[:=]\s*["']?([^\s"']+)["']?`)
)

// SanitizeText scrubs sensitive credentials, tokens, and private keys from prompts
func SanitizeText(input string, customPatterns []string) string {
	if input == "" {
		return ""
	}

	sanitized := privateKeyRegex.ReplaceAllString(input, "[REDACTED_PRIVATE_KEY]")
	sanitized = awsAccessKeyRegex.ReplaceAllString(sanitized, "[REDACTED_AWS_KEY]")
	sanitized = githubTokenRegex.ReplaceAllString(sanitized, "[REDACTED_GITHUB_TOKEN]")
	sanitized = jwtRegex.ReplaceAllString(sanitized, "[REDACTED_JWT_TOKEN]")
	sanitized = bearerRegex.ReplaceAllString(sanitized, "Bearer [REDACTED_TOKEN]")
	sanitized = urlPassRegex.ReplaceAllString(sanitized, "$1:[REDACTED_PASSWORD]@")
	sanitized = envSecretRegex.ReplaceAllString(sanitized, "$1=[REDACTED]")

	for _, pat := range customPatterns {
		if pat == "" {
			continue
		}
		if re, err := regexp.Compile(pat); err == nil {
			sanitized = re.ReplaceAllString(sanitized, "[REDACTED]")
		}
	}

	return sanitized
}
