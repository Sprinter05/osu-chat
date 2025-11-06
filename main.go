package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sprinter05/osu-chat/api"
)

func main() {
	config := getConfig()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
			MaxIdleConns:      5,
			IdleConnTimeout:   30 * time.Second,
		},
		Timeout: time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			print(req, via)
			return nil
		},
	}

	token, err := api.RequestToken(client, config.OAuth)
	if err != nil {
		log.Fatal(err)
	}

	err = api.DeleteToken(client, token)
	if err != nil {
		log.Fatal(err)
	}
}
