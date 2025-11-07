package main

import (
	"time"

	"github.com/Sprinter05/osu-chat/api"
)

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

// Removes the duration on the token to be used
// in API requests
func configToToken(token Token) api.Token {
	return api.Token{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}
