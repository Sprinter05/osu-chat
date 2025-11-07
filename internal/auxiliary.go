package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

/* ERROR HANDLING */

func HTTPError(res *http.Response) error {
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

/* TIME FORMATING ACCORDING TO ISO 8601 */

const TIME_FORMAT string = "2006-01-02T15:04:05Z0700"

func ParseTime(s string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, s)
}

func FormatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

/* HTTP CLIENTS */

func DefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
			MaxIdleConns:      5,
			IdleConnTimeout:   30 * time.Second,
		},
		Timeout: time.Minute,
	}
}

/* HTTP HEADERS */

func SetContentHeaders(hd *http.Header) {
	hd.Add("Accept", "application/json")
	hd.Add("Content-Type", "application/json")
}
