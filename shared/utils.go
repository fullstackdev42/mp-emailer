package shared

import "time"

// ParseDateTime parses a date string
func ParseDateTime(dateStr string) (time.Time, error) {
	if dateStr == "" || dateStr == "0000-00-00 00:00:00" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05", dateStr)
}
