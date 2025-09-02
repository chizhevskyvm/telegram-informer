package utils

import (
	"strings"
	"time"
)

func FromTimeZone(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.TrimSpace(s))
}

func ParseDateTz(s, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	return time.ParseInLocation("2006-01-02", strings.TrimSpace(s), loc)
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatTime(t time.Time) string {
	return t.Format("15:04")
}
