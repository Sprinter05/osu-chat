package main

import (
	"flag"
	"log"
	"time"

	"github.com/Sprinter05/osu-chat/api"
	"github.com/Sprinter05/osu-chat/internal"
)

/* SETUP */

var configFile string

func init() {
	flag.StringVar(
		&configFile, "config", "config.json",
		"Configuration file to load, must be in JSON format.",
	)
}

func login(config *Config) (api.Token, error) {
	if config.Token != nil {
		// TODO: refresh token if necessary
		return configToToken(*config.Token), nil
	}

	token, err := api.RetrieveToken(config.OAuth)
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

	err = internal.SaveConfig(configFile, *config)
	if err != nil {
		log.Fatal(err)
	}

	return configToToken(*config.Token), nil
}

func main() {
	config := new(Config)
	client := internal.DefaultClient()
	err := internal.GetConfig(configFile, config)
	if err != nil {
		log.Fatal(err)
	}

	token, err := login(config)
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
