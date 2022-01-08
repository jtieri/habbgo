package date

import (
	"fmt"
	"time"
)

// GetCurrentDateTime returns the current date and time formatted as dd-MM-yyyy HH:mm:ss
func GetCurrentDateTime() string {
	t := time.Now()
	formatted := fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(), t.Second())
	return formatted
}

// GetCurrentDate returns the current date formatted as dd-MM-yyyy
func GetCurrentDate() string {
	t := time.Now()
	formatted := fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
	return formatted
}
