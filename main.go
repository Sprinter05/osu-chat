package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sprinter05/osu-chat/api"
)

// TODO: move to dedicated OAuth serer (client secret is exposed)

func defaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
			MaxIdleConns:      5,
			IdleConnTimeout:   30 * time.Second,
		},
		Timeout: time.Minute,
	}
}

func login(client *http.Client, config *Config) (api.Token, error) {
	if config.Token != nil {
		// TODO: refresh token if necessary
		return configToToken(*config.Token), nil
	}

	token, err := api.RequestToken(client, config.OAuth)
	if err != nil {
		log.Fatal(err)
	}

	expiration := time.Now().Add(
		time.Duration(token.ExpiresIn) * time.Second,
	)

	config.Token = &Token{
		TokenType:      token.TokenType,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		ExpirationDate: expiration,
	}

	err = saveConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	return configToToken(*config.Token), nil
}

func main() {
	client := defaultClient()
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	token, err := login(client, &config)
	if err != nil {
		log.Fatal(err)
	}

	list, err := api.GetChannelList(client, token)
	if err != nil {
		log.Fatal(err)
	}
	print(list)

	err = api.DeleteToken(client, token)
	if err != nil {
		log.Fatal(err)
	}
}
