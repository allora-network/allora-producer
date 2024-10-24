package util

import (
	"runtime/debug"
)

// readBuildInfo is a variable that references debug.ReadBuildInfo.
// It can be overridden for testing purposes.
var readBuildInfo = debug.ReadBuildInfo

// Revision retrieves the VCS revision from build info.
// Available only if the build was done with -buildvcs flag.
func Revision() string {
	if info, ok := readBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}
