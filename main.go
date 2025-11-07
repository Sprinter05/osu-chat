package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Sprinter05/osu-chat/api"
	"github.com/Sprinter05/osu-chat/callback/oauth"
	"github.com/Sprinter05/osu-chat/internal"
)

/* SETUP */

// TODO: Docker for OAuth intermediate server

var configFile string

func init() {
	flag.StringVar(
		&configFile, "config", "config.json",
		"Configuration file to load, must be in JSON format.",
	)
	flag.Parse()
}

/* LOGIN */

// Returns the API Token to be used in petitions
// and saves the oauth token in the configuration file
func handleToken(config Config, token oauth.Token) (api.Token, error) {
	expiration := time.Now().Add(
		time.Duration(token.ExpiresIn) * time.Second,
	)

	// Token for the configuration file
	toConfig := Token{
		TokenType:      token.TokenType,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		ExpirationDate: expiration,
	}
	config.Token = &toConfig

	// Save token
	err := internal.SaveConfig(configFile, config)
	if err != nil {
		return api.Token{}, err
	}

	// Return API Token
	return configToAPIToken(toConfig), nil
}

// Login into the osu API, authorizing and refreshing if necessary
func login(cl *http.Client, config *Config) (api.Token, error) {
	// Token alerady exists
	if config.Token != nil {
		// Token has expired
		if time.Now().After(config.Token.ExpirationDate) {
			convert := configToAPIToken(*config.Token)

			// Refresh a new token
			newToken, err := api.RefreshToken(cl, config.OAuth, convert)
			if err != nil {
				return api.Token{}, err
			}

			// Save to configuration and return API token
			return handleToken(*config, newToken)
		}

		// Just return API token since it has not expired
		return configToAPIToken(*config.Token), nil
	}

	// There is no token so we retrieve it
	token, err := api.RetrieveToken(config.OAuth)
	if err != nil {
		log.Fatal(err)
	}

	// Save new token to configuration file and return API token
	return handleToken(*config, token)
}

func main() {
	config := new(Config)
	client := internal.DefaultClient()
	err := internal.GetConfig(configFile, config)
	if err != nil {
		log.Fatal(err)
	}

	token, err := login(client, config)
	if err != nil {
		log.Fatal(err)
	}

	list, err := api.GetChannelList(client, token)
	if err != nil {
		log.Fatal(err)
	}
	print(list)
}
