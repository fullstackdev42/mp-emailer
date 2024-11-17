package version

// Info holds version information for the application
type Info struct {
	Version   string
	BuildDate string
	Commit    string
}

// Status returns version information as a map
func (i Info) Status() map[string]string {
	return map[string]string{
		"version":   i.Version,
		"buildDate": i.BuildDate,
		"commit":    i.Commit,
	}
}
