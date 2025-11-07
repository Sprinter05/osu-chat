package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Errors

func httpErr(res *http.Response) error {
	msg, _ := io.ReadAll(res.Body)
	if len(msg) != 0 {
		m := make(map[string]string)
		err := json.Unmarshal(msg, &m)
		if err == nil {
			v, ok := m["error"]
			if ok {
				return fmt.Errorf(
					"http returned %s: %s",
					res.Status, v,
				)
			}
		}
	}

	return errors.New(res.Status)
}

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
		"%s %s",
		token.TokenType,
		token.AccessToken,
	))
}
