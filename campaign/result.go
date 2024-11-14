package campaign

// Result represents the interface for database query results
type Result interface {
	// Scan copies the result into the provided destination
	Scan(dest interface{}) Result
	// Error returns any error that occurred during the query
	Error() error
}

// DatabaseResult implements the Result interface
type DatabaseResult struct {
	err error
}

// Scan implements Result.Scan
func (r *DatabaseResult) Scan(_ interface{}) Result {
	// Implementation would depend on your actual database
	return r
}

// Error implements Result.Error
func (r *DatabaseResult) Error() error {
	return r.err
}
