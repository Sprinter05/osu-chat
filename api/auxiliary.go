package api

import (
	"fmt"
	"net/http"
	"time"
)

// Time according to ISO 8601

const TIME_FORMAT string = "2006-01-02T15:04:05Z0700"

func parseTime(s string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, s)
}

func formatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

// Authentication

func setGenericHeaders(hd *http.Header, token Token) {
	hd.Set("Content-Type", "application/json")
	hd.Set("Accept", "application/json")
	hd.Set("Authorization", fmt.Sprintf(
		"Bearer %s",
		token.AccessToken,
	))
}
