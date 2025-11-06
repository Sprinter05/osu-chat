package main

import (
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
			MaxIdleConns:      5,
			IdleConnTimeout:   30 * time.Second,
		},
		Timeout: time.Minute,
	}
}
