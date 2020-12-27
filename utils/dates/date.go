package dates

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDBLayout   = "2006-01-02 15:04:05"
)

// GetNow returns current UTC time
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString return a formated string of current UTC time
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

// GetNowDBFormat returns a DB formated string of current UTC time
func GetNowDBFormat() string {
	return GetNow().Format(apiDBLayout)
}
