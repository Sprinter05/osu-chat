package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sprinter05/osu-chat/api"
)

func login(client *http.Client, config *Config) error {
	if config.Token != nil {
		// TODO: refresh token if necessary
		return nil
	}

	token, err := api.RequestToken(client, config.OAuth)
	if err != nil {
		log.Fatal(err)
	}

	config.Token = &token
	err = saveConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	err = login(client, &config)
	if err != nil {
		log.Fatal(err)
	}

	err = api.DeleteToken(client, *config.Token)
	if err != nil {
		log.Fatal(err)
	}
}
