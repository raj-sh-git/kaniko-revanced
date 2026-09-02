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
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const autoHealSystemPrompt = "You are the Autonomous Auto-Healing Engine for kaniko-revanced.\n" +
	"A container build step has failed. Your task is to fix the Dockerfile so that the build succeeds.\n\n" +
	"Rules:\n" +
	"1. Fix the root cause (e.g. missing package manager dependencies, wrong flag, missing build toolchain, syntax error, missing directory creation).\n" +
	"2. Maintain all original intent and functionality of the image.\n" +
	"3. Respond in this EXACT format:\n\n" +
	"EXPLANATION: <One concise sentence explaining the fix applied>\n\n" +
	"```dockerfile\n" +
	"<The complete, corrected Dockerfile with no missing stages or truncation>\n" +
	"```\n"

var dockerfileBlockRegex = regexp.MustCompile("(?s)```(?:dockerfile|docker|sh)?\\s*\\n(.*?)\\n```")

// RunAutoHeal requests an auto-healed Dockerfile from the LLM and computes the unified diff
func RunAutoHeal(ctx context.Context, client *Client, errCtx BuildErrorContext) (*AutoHealResult, error) {
	var userPrompt strings.Builder
	userPrompt.WriteString("The following Dockerfile failed to build in kaniko-revanced. Please fix it.\n\n")

	if errCtx.FailedStageName != "" {
		userPrompt.WriteString(fmt.Sprintf("Failed Stage: %s (Index %d)\n", errCtx.FailedStageName, errCtx.FailedStageIndex))
	}
	if errCtx.FailedCommand != "" {
		userPrompt.WriteString(fmt.Sprintf("Failed Instruction: `%s`\n", errCtx.FailedCommand))
	}
	if errCtx.ErrorMessage != "" {
		userPrompt.WriteString(fmt.Sprintf("Error Message: %s\n", errCtx.ErrorMessage))
	}
	if errCtx.RecentLogs != "" {
		userPrompt.WriteString(fmt.Sprintf("Logs/Stderr:\n```\n%s\n```\n\n", errCtx.RecentLogs))
	}

	userPrompt.WriteString(fmt.Sprintf("Original Dockerfile:\n```dockerfile\n%s\n```\n", errCtx.DockerfileContent))

	rawResponse, err := client.Complete(ctx, autoHealSystemPrompt, userPrompt.String())
	if err != nil {
		return nil, errors.Wrap(err, "auto-heal LLM request failed")
	}

	explanation, patchedDockerfile, err := parseAutoHealResponse(rawResponse)
	if err != nil {
		return nil, err
	}

	diff := generateUnifiedDiff(errCtx.DockerfileContent, patchedDockerfile)

	return &AutoHealResult{
		PatchedDockerfile: patchedDockerfile,
		Explanation:       explanation,
		Diff:              diff,
	}, nil
}

func parseAutoHealResponse(response string) (string, string, error) {
	explanation := "Applied automated build fix"
	lines := strings.Split(response, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(strings.ToUpper(trimmed), "EXPLANATION:") {
			explanation = strings.TrimSpace(trimmed[len("EXPLANATION:"):])
			break
		}
	}

	matches := dockerfileBlockRegex.FindStringSubmatch(response)
	if len(matches) < 2 {
		return "", "", errors.New("failed to extract corrected Dockerfile code block from LLM response")
	}

	patchedDockerfile := strings.TrimSpace(matches[1])
	if patchedDockerfile == "" {
		return "", "", errors.New("extracted Dockerfile is empty")
	}

	return explanation, patchedDockerfile, nil
}

func generateUnifiedDiff(original, modified string) string {
	origLines := strings.Split(original, "\n")
	modLines := strings.Split(modified, "\n")

	var diff strings.Builder
	diff.WriteString("--- Original Dockerfile\n+++ Patched Dockerfile (AI Auto-Healed)\n")

	// Line by line comparison
	maxLen := len(origLines)
	if len(modLines) > maxLen {
		maxLen = len(modLines)
	}

	for i := 0; i < maxLen; i++ {
		var origLine, modLine string
		hasOrig := i < len(origLines)
		hasMod := i < len(modLines)

		if hasOrig {
			origLine = origLines[i]
		}
		if hasMod {
			modLine = modLines[i]
		}

		if hasOrig && hasMod && origLine == modLine {
			diff.WriteString(fmt.Sprintf("  %s\n", origLine))
		} else {
			if hasOrig {
				diff.WriteString(fmt.Sprintf("- %s\n", origLine))
			}
			if hasMod {
				diff.WriteString(fmt.Sprintf("+ %s\n", modLine))
			}
		}
	}

	return diff.String()
}
