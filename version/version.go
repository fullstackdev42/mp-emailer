package version

var (
	// current holds the current version information
	current = Info{
		Version:   "dev",
		BuildDate: "unknown",
		Commit:    "none",
	}
)

// Get returns the current version information
func Get() Info {
	return current
}

// Set updates the version information
func Set(version, buildDate, commit string) {
	if version != "" {
		current.Version = version
	}
	if buildDate != "" {
		current.BuildDate = buildDate
	}
	if commit != "" {
		current.Commit = commit
	}
}
