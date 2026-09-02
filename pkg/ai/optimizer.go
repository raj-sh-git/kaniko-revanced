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
)

const lintSystemPrompt = `You are the Pre-Flight Dockerfile Linter and Optimizer for kaniko-revanced.
Review the provided Dockerfile before build execution and highlight:
1. ⚠️ Critical Anti-Patterns (hardcoded secrets in ENV, running everything as root, insecure downloads without checksums)
2. ⚡ Cache Inefficiencies (copying dynamic files before dependency lockfiles like package.json or go.mod)
3. 📉 Layer & Size Waste (unnecessary package manager caches, missing multi-stage builds)
4. 💡 Quick Wins (immediate one-line changes)

Keep the analysis concise, actionable, and formatted with clear emoji headings.`

// RunLint executes pre-flight static analysis on a Dockerfile
func RunLint(ctx context.Context, client *Client, dockerfileContent, targetArch string) (string, error) {
	userPrompt := fmt.Sprintf("Perform pre-flight static analysis on this Dockerfile (Target Arch: %s):\n\n```dockerfile\n%s\n```", targetArch, dockerfileContent)
	return client.Complete(ctx, lintSystemPrompt, userPrompt)
}
