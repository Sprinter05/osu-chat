package wrappers

import (
	"log"
	"net/http"
	"time"

	"github.com/Sprinter05/osu-chat/api"
	"github.com/Sprinter05/osu-chat/internal/conf"
)

/* LOGIN */

// Login into the osu API, authorizing and refreshing if necessary
func Login(cl *http.Client, configFile string, config *conf.Config) (api.Token, error) {
	// Token alerady exists
	if config.Token != nil {
		// Token has expired
		if time.Now().After(config.Token.ExpirationDate) {
			convert := conf.ConfigToAPIToken(*config.Token)

			// Refresh a new token
			newToken, err := api.RefreshToken(cl, config.OAuth, convert)
			if err != nil {
				return api.Token{}, err
			}

			// Save to configuration and return API token
			return conf.HandleToken(configFile, *config, newToken)
		}

		// Just return API token since it has not expired
		return conf.ConfigToAPIToken(*config.Token), nil
	}

	// There is no token so we retrieve it
	token, err := api.RetrieveToken(config.OAuth)
	if err != nil {
		log.Fatal(err)
	}

	// Save new token to configuration file and return API token
	return conf.HandleToken(configFile, *config, token)
}
