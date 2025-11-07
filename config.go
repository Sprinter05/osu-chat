package main

import (
	"time"

	"github.com/Sprinter05/osu-chat/api"
)

// TODO: Add arguments
const CONFIG_FILE = "config/config.json"

type Token struct {
	TokenType      string    `json:"token_type"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type Config struct {
	OAuth api.OAuth `json:"api"`
	Token *Token    `json:"token,omitempty"`
}

func configToToken(token Token) api.Token {
	return api.Token{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}
