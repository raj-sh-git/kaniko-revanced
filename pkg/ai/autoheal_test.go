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

func TestParseAutoHealResponse(t *testing.T) {
	rawResponse := `Here is the fix:

EXPLANATION: Installed build-base and libffi-dev packages required for compilation.

` + "```dockerfile" + `
FROM alpine:3.20
RUN apk add --no-cache build-base libffi-dev python3
RUN pip install cryptography --break-system-packages
` + "```"

	explanation, patched, err := parseAutoHealResponse(rawResponse)
	if err != nil {
		t.Fatalf("unexpected error parsing response: %v", err)
	}

	if !strings.Contains(explanation, "Installed build-base") {
		t.Errorf("unexpected explanation: %q", explanation)
	}

	if !strings.Contains(patched, "RUN apk add --no-cache build-base libffi-dev python3") {
		t.Errorf("unexpected patched Dockerfile: %q", patched)
	}
}

func TestGenerateUnifiedDiff(t *testing.T) {
	orig := "FROM alpine:3.20\nRUN pip install cryptography"
	mod := "FROM alpine:3.20\nRUN apk add --no-cache gcc musl-dev\nRUN pip install cryptography"

	diff := generateUnifiedDiff(orig, mod)
	if !strings.Contains(diff, "+ RUN apk add --no-cache gcc musl-dev") {
		t.Errorf("diff did not contain added line: %s", diff)
	}
}
