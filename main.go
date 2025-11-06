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
	}

	list, err := api.GetChannelList(client, config.OAuth)
	if err != nil {
		log.Fatal(err)
	}

	print(list)
}
