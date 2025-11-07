package conf

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Sprinter05/osu-chat/api"
	"github.com/Sprinter05/osu-chat/callback/oauth"
)

const DefaultPerms = 0755

type Config struct {
	OAuth api.OAuth `json:"api"`
	Token *Token    `json:"token,omitempty"`
}

/* TOKENS */

// Token type to be stored in the configuration
type Token struct {
	TokenType      string    `json:"token_type"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpirationDate time.Time `json:"expiration_date"`
}

// Removes the duration on the token to be used
// in API requests
func ConfigToAPIToken(token Token) api.Token {
	return api.Token{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

// Returns the API Token to be used in petitions
// and saves the oauth token in the configuration file
func HandleToken(file string, config Config, token oauth.Token) (api.Token, error) {
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
	err := SaveConfig(file, config)
	if err != nil {
		return api.Token{}, err
	}

	// Return API Token
	return ConfigToAPIToken(toConfig), nil
}

/* CONFIG FILES */

// Gets the config into the "config" variable.
// Said variable must be a pointer
func GetConfig(path string, config any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(config)
	return nil
}

// Saves the config into the specified file.
// The provided configuration must not be a pointer
func SaveConfig(path string, config any) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, DefaultPerms)
	return err
}
