package main

import "runtime/debug"

// shortRevision is how many characters of the commit hash version shows.
const shortRevision = 7

// version returns "vitae <commit>" from the build info, or "vitae dev".
func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= shortRevision {
				return "vitae " + s.Value[:shortRevision]
			}
		}
	}
	return "vitae dev"
}
