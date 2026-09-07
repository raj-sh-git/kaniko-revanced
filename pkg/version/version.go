/*
Copyright 2018 Google LLC

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

package version

import (
	"fmt"
	"runtime"
)

// ProjectName is the name of the fork
const ProjectName = "kaniko-revanced"

// Set with LDFLAGS
var (
	version   = "v0.2.0"
	commit    = "unknown"
	buildDate = "unknown"
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func BuildDate() string {
	return buildDate
}

func FullVersion() string {
	return fmt.Sprintf("%s %s (commit: %s, built: %s, %s/%s, go: %s)",
		ProjectName,
		version,
		commit,
		buildDate,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version(),
	)
}
