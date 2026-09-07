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
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiBlue   = "\033[34m"
	ansiCyan   = "\033[36m"
)

func isGitLabCI() bool {
	return os.Getenv("GITLAB_CI") == "true" || os.Getenv("CI_SERVER") == "yes"
}

func isGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// PrintHeader renders a stylish banner or CI collapsible group
func PrintHeader(title string) {
	if isGitLabCI() {
		secName := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
		fmt.Printf("\033[0Ksection_start:%d:%s[collapsed=true]\r\033[0K%s%s%s%s\n", time.Now().Unix(), secName, ansiCyan, ansiBold, title, ansiReset)
		return
	}

	if isGitHubActions() {
		fmt.Printf("::group::%s\n", title)
		return
	}

	width := 72
	fmt.Println()
	fmt.Printf("%s%s╔%s╗%s\n", ansiCyan, ansiBold, strings.Repeat("═", width-2), ansiReset)
	padding := (width - 4 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s║ %s%s%s ║%s\n", ansiCyan, ansiBold, strings.Repeat(" ", padding), title, strings.Repeat(" ", width-4-len(title)-padding), ansiReset)
	fmt.Printf("%s%s╚%s╝%s\n", ansiCyan, ansiBold, strings.Repeat("═", width-2), ansiReset)
}

// PrintFooter closes collapsible CI groups if open
func PrintFooter(title string) {
	if isGitLabCI() {
		secName := strings.ToLower(strings.ReplaceAll(title, " ", "_"))
		fmt.Printf("\033[0Ksection_end:%d:%s\r\033[0K\n", time.Now().Unix(), secName)
		return
	}
	if isGitHubActions() {
		fmt.Println("::endgroup::")
	}
}

// PrintDiff renders a colorized unified diff
func PrintDiff(diff string) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			fmt.Printf("%s%s%s%s\n", ansiCyan, ansiBold, line, ansiReset)
		} else if strings.HasPrefix(line, "-") {
			fmt.Printf("%s%s%s\n", ansiRed, line, ansiReset)
		} else if strings.HasPrefix(line, "+") {
			fmt.Printf("%s%s%s\n", ansiGreen, line, ansiReset)
		} else {
			fmt.Printf("  %s\n", line)
		}
	}
}

// PrintReport renders AI markdown report
func PrintReport(content string) {
	fmt.Println()
	fmt.Println(content)
	fmt.Println()
}
