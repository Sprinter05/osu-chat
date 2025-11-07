package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

/* ERROR HANDLING */

// Returns an HTTP error as a go error
// with more information if available
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

/* TIME FORMATING */

// Format according to ISO 8601
const TimeFormat string = "2006-01-02T15:04:05Z0700"

func ParseTime(s string) (time.Time, error) {
	return time.Parse(TimeFormat, s)
}

func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

/* HTTP  */

// Check if a port is in use by dialing that port in localhost
// and checking if an answer is received
func CheckPortInUse(port uint16) bool {
	portStr := strconv.FormatInt(int64(port), 10)
	address := net.JoinHostPort("localhost", portStr)
	conn, err := net.DialTimeout("tcp4", address, time.Second)
	if err != nil {
		// Port has nothing running on it
		return true
	}

	if conn == nil {
		// Connection is null meaning it doesn't exist
		return true
	}

	// We close the connection if it was obtained
	// and return that a connection was found (not available)
	conn.Close()
	return false
}

// Returns a default http client for requests
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

// Sets the headers to accept and send "application/json"
func SetContentHeaders(hd *http.Header) {
	hd.Add("Accept", "application/json")
	hd.Add("Content-Type", "application/json")
}
