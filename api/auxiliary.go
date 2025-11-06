package api

import "time"

// Time according to ISO 8601

func ParseTime(s string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, s)
}

func FormatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}
