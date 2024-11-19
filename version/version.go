package version

// NewInfo creates and returns a default version Info
func NewInfo() Info {
	return Info{
		Version:   "dev",
		BuildDate: "unknown",
		Commit:    "none",
	}
}

// Get returns the current version information
func Get() Info {
	return NewInfo()
}

// Set returns a new Info with updated fields
func Set(version, buildDate, commit string) Info {
	info := NewInfo()
	if version != "" {
		info.Version = version
	}
	if buildDate != "" {
		info.BuildDate = buildDate
	}
	if commit != "" {
		info.Commit = commit
	}
	return info
}
