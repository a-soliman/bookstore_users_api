package dates

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
)

// GetNow returns current UTC time
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString return a formated string of current UTC time
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}
