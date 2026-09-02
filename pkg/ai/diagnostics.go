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
	"strings"
)

const diagnosticSystemPrompt = `You are the AI Container Diagnostic and Optimization Engine for kaniko-revanced (a fast, daemonless container builder).
Your job is to analyze container build failures and provide:
1. ROOT CAUSE: Clear, concise explanation of why the build failed.
2. SUGGESTED FIX: Specific commands or configuration changes to fix the error.
3. DOCKERFILE DIFF: A clean unified diff showing the exact changes to apply to the Dockerfile.
4. IMAGE SIZE REDUCTION: Actionable recommendations to reduce final container image size (e.g., multi-stage builds, removing cache/temp files, using --no-cache, stripping debug symbols).
5. CODEBASE & DOCKERFILE IMPROVEMENTS: Best practices for caching, security (non-root USER, secret mounts), and layer efficiency.

Format your response clearly using markdown with clear headings:
## 🔴 Build Failure Root Cause
## 💡 Recommended Fix & Dockerfile Diff
## 📉 Image Size Reduction Suggestions
## ⭐ Codebase & Best Practice Recommendations`

// RunDiagnostics executes full failure diagnosis and optimization review
func RunDiagnostics(ctx context.Context, client *Client, errCtx BuildErrorContext) (string, error) {
	var userPrompt strings.Builder
	userPrompt.WriteString("A container build has failed in kaniko-revanced.\n\n")

	userPrompt.WriteString("### Build Context:\n")
	if errCtx.DockerfilePath != "" {
		userPrompt.WriteString(fmt.Sprintf("- Dockerfile Path: %s\n", errCtx.DockerfilePath))
	}
	if errCtx.FailedStageName != "" {
		userPrompt.WriteString(fmt.Sprintf("- Failed Stage: %s (Index: %d)\n", errCtx.FailedStageName, errCtx.FailedStageIndex))
	}
	if errCtx.TargetOS != "" || errCtx.TargetArch != "" {
		userPrompt.WriteString(fmt.Sprintf("- Platform: %s/%s\n", errCtx.TargetOS, errCtx.TargetArch))
	}
	if errCtx.FailedCommand != "" {
		userPrompt.WriteString(fmt.Sprintf("- Failed Instruction: `%s`\n", errCtx.FailedCommand))
	}
	if errCtx.ErrorMessage != "" {
		userPrompt.WriteString(fmt.Sprintf("- Error Message: %s\n", errCtx.ErrorMessage))
	}

	if errCtx.RecentLogs != "" {
		userPrompt.WriteString("\n### Execution Logs / Stderr (Tail):\n```\n")
		userPrompt.WriteString(errCtx.RecentLogs)
		userPrompt.WriteString("\n```\n")
	}

	if errCtx.DockerfileContent != "" {
		userPrompt.WriteString("\n### Current Dockerfile:\n```dockerfile\n")
		userPrompt.WriteString(errCtx.DockerfileContent)
		userPrompt.WriteString("\n```\n")
	}

	userPrompt.WriteString("\nPlease provide the complete root-cause diagnostic, fix diff, image size reduction recommendations, and Dockerfile quality advice.")

	return client.Complete(ctx, diagnosticSystemPrompt, userPrompt.String())
}

// RunSuccessDiagnostics reviews a successfully built Dockerfile for size reduction & best practices
func RunSuccessDiagnostics(ctx context.Context, client *Client, dockerfileContent, targetArch string) (string, error) {
	systemPrompt := `You are the AI Container Optimization Engine for kaniko-revanced.
Analyze this Dockerfile and provide actionable recommendations for:
1. 📉 Image Size Reduction (multi-stage builds, cache purging, smaller base images)
2. ⚡ Build Cache Optimization (layer ordering)
3. 🔒 Security Best Practices (non-root users, secret mounts)`

	userPrompt := fmt.Sprintf("Analyze this successfully built Dockerfile for optimization opportunities (Platform: %s):\n\n```dockerfile\n%s\n```", targetArch, dockerfileContent)

	return client.Complete(ctx, systemPrompt, userPrompt)
}
