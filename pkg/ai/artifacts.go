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
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ExportArtifacts writes out the patched Dockerfile, report, patch, and audit JSON to artifactDir
func ExportArtifacts(artifactDir string, patchedDockerfile, reportContent, unifiedDiff string, auditLogs []AIAuditEntry) error {
	if artifactDir == "" {
		return nil
	}

	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create artifact directory %s", artifactDir)
	}

	if patchedDockerfile != "" {
		pPath := filepath.Join(artifactDir, "Dockerfile.fixed")
		if err := os.WriteFile(pPath, []byte(patchedDockerfile), 0644); err != nil {
			logrus.Warnf("Failed to write %s: %v", pPath, err)
		} else {
			logrus.Infof("Saved patched Dockerfile to %s", pPath)
		}
	}

	if reportContent != "" {
		rPath := filepath.Join(artifactDir, "optimization_report.md")
		if err := os.WriteFile(rPath, []byte(reportContent), 0644); err != nil {
			logrus.Warnf("Failed to write %s: %v", rPath, err)
		} else {
			logrus.Infof("Saved AI diagnostic report to %s", rPath)
		}
	}

	if unifiedDiff != "" {
		dPath := filepath.Join(artifactDir, "changes.patch")
		if err := os.WriteFile(dPath, []byte(unifiedDiff), 0644); err != nil {
			logrus.Warnf("Failed to write %s: %v", dPath, err)
		} else {
			logrus.Infof("Saved unified diff patch to %s", dPath)
		}
	}

	if len(auditLogs) > 0 {
		aPath := filepath.Join(artifactDir, "ai-audit.json")
		if data, err := json.MarshalIndent(auditLogs, "", "  "); err == nil {
			if err := os.WriteFile(aPath, data, 0644); err != nil {
				logrus.Warnf("Failed to write %s: %v", aPath, err)
			} else {
				logrus.Infof("Saved AI prompt audit history to %s", aPath)
			}
		}
	}

	return nil
}
