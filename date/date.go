package date

import (
	"time"
)

const (
	dateTimeFormat = "02-01-2006 15:04:05"
	dateFormat     = "02-01-2006"
)

// CurrentDateTime returns the current date and time formatted as dd-MM-yyyy HH:mm:ss.
func CurrentDateTime() string {
	t := time.Now()
	return t.Format(dateTimeFormat)
}

// CurrentDate returns the current date formatted as dd-MM-yyyy.
func CurrentDate() string {
	t := time.Now()
	return t.Format(dateFormat)
}
