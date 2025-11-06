package api

import (
	"fmt"
	"net/http"
	"time"
)

// Time according to ISO 8601

func parseTime(s string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, s)
}

func formatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

// Authentication

func setGenericHeaders(hd *http.Header, oauth OAuth) {
	hd.Set("Content-Type", "application/json")
	hd.Set("Accept", "application/json")
	hd.Set("Authorization", fmt.Sprintf(
		"Bearer %s",
		oauth.TokenSecret,
	))
}
