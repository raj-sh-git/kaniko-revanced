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

	"github.com/pkg/errors"
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

const diagnoseApplySystemPrompt = "You are the AI Dockerfile Optimization Engine for kaniko-revanced.\n" +
	"Your task is to analyze the provided Dockerfile and rewrite it to be fully optimized for:\n" +
	"1. Minimal image size (multi-stage builds, --no-cache, purging package caches like /var/cache/apk/* or /var/lib/apt/lists/*).\n" +
	"2. Layer caching efficiency (copying lockfiles first).\n" +
	"3. Security best practices (non-root USER).\n\n" +
	"Respond in this EXACT format:\n\n" +
	"EXPLANATION: <One concise sentence summarizing the optimizations applied>\n\n" +
	"```dockerfile\n" +
	"<The complete, optimized Dockerfile>\n" +
	"```\n"

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

	if errCtx.PrivacyMode == "strict" {
		userPrompt.WriteString("\n(Privacy Mode Active: Full Dockerfile content omitted from prompt)\n")
	} else if errCtx.DockerfileContent != "" {
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

// RunDiagnoseAndApply optimizes the Dockerfile and returns the patched contents + unified diff
func RunDiagnoseAndApply(ctx context.Context, client *Client, dockerfileContent, targetArch string) (*AutoHealResult, error) {
	userPrompt := fmt.Sprintf("Optimize this Dockerfile for minimal image size, layer caching, and security (Platform: %s):\n\n```dockerfile\n%s\n```", targetArch, dockerfileContent)

	rawResponse, err := client.Complete(ctx, diagnoseApplySystemPrompt, userPrompt)
	if err != nil {
		return nil, errors.Wrap(err, "diagnose-apply LLM request failed")
	}

	explanation, optimizedDockerfile, err := parseAutoHealResponse(rawResponse)
	if err != nil {
		return nil, err
	}

	if err := ValidatePatchedDockerfile(dockerfileContent, optimizedDockerfile); err != nil {
		return nil, errors.Wrap(err, "diagnose-apply validation failed")
	}

	diff := generateUnifiedDiff(dockerfileContent, optimizedDockerfile)

	return &AutoHealResult{
		PatchedDockerfile: optimizedDockerfile,
		Explanation:       explanation,
		Diff:              diff,
	}, nil
}
