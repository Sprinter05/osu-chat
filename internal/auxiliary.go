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

/* HTTP  */

func CheckPortInUse(port uint16) bool {
	portStr := strconv.FormatInt(int64(port), 10)
	address := net.JoinHostPort("localhost", portStr)
	conn, err := net.DialTimeout("tcp4", address, time.Second)
	if err != nil {
		return true
	}

	if conn == nil {
		return true
	}

	conn.Close()
	return false
}

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

func SetContentHeaders(hd *http.Header) {
	hd.Add("Accept", "application/json")
	hd.Add("Content-Type", "application/json")
}
