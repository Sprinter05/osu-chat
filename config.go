package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Sprinter05/osu-chat/api"
)

// TODO: Add arguments
const CONFIG_FILE = "config/config.json"
const DEFAULT_PERMS = 0755

type Token struct {
	TokenType      string    `json:"token_type"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type Config struct {
	OAuth api.OAuth `json:"oauth"`
	Token *Token    `json:"token,omitempty"`
}

func configToToken(token Token) api.Token {
	return api.Token{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

func getConfig() (config Config, err error) {
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		return config, err
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&config)
	return config, nil
}

func saveConfig(config Config) error {
	data, err := json.MarshalIndent(config, "\t", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(CONFIG_FILE, data, DEFAULT_PERMS)
	return err
}
