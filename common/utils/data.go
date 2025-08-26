package utils

import (
	"strings"
	"time"
)

func ParseDateLocal(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", strings.TrimSpace(s), time.Local)
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
