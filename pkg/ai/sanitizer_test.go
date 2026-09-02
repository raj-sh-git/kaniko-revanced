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
	"strings"
	"testing"
)

func TestSanitizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		patterns []string
		mustNot  string
		mustHave string
	}{
		{
			name:     "AWS Access Key",
			input:    "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE",
			mustNot:  "AKIAIOSFODNN7EXAMPLE",
			mustHave: "[REDACTED_AWS_KEY]",
		},
		{
			name:     "GitHub Token",
			input:    "RUN git clone https://ghp_123456789012345678901234567890123456@github.com/org/repo.git",
			mustNot:  "ghp_123456789012345678901234567890123456",
			mustHave: "[REDACTED_GITHUB_TOKEN]",
		},
		{
			name:     "Password in URL",
			input:    "RUN curl https://admin:SuperSecretPass123@internal.corp.net/data.tar.gz",
			mustNot:  "SuperSecretPass123",
			mustHave: "[REDACTED_PASSWORD]",
		},
		{
			name:     "Custom Redaction Pattern",
			input:    "CONFIDENTIAL_CUSTOMER_UUID_ABC999",
			patterns: []string{`UUID_[A-Z0-9]+`},
			mustNot:  "UUID_ABC999",
			mustHave: "[REDACTED]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeText(tc.input, tc.patterns)
			if strings.Contains(got, tc.mustNot) {
				t.Errorf("SanitizeText(%q) did not redact secret %q, got %q", tc.input, tc.mustNot, got)
			}
			if !strings.Contains(got, tc.mustHave) {
				t.Errorf("SanitizeText(%q) did not contain expected placeholder %q, got %q", tc.input, tc.mustHave, got)
			}
		})
	}
}
