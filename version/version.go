// SPDX-License-Identifier: Apache-2.0

package version

import "strings"

// These properties are filled through -ldflags upon building.
// See the Makefile.

// GitSummary is a long, descriptive git version like 0.3.1 or 0.3.1-12-g3257b77
var GitSummary string

// BuildDate in ISO8601 ; looks like this 2021-10-20T12:24:58Z
var BuildDate string

// GetVersion returns the fully described git version of this bot.
func GetVersion() string {
	if GitSummary == "" {
		return "N/A"
	}
	return strings.ReplaceAll(GitSummary, "dirty", "edge")
}
