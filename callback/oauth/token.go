package oauth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sprinter05/osu-chat/internal"
)

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Request a token to the API based on the standard OAuth 2.0 specification
func requestToken(cl *http.Client, config Config, code string) (Token, error) {
	values := map[string]string{
		"client_id":     strconv.FormatInt(int64(config.OAuth.ClientId), 10),
		"client_secret": config.OAuth.TokenSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  config.CallbackURL,
	}

	body, err := json.Marshal(values)
	if err != nil {
		return Token{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		config.TokenURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return Token{}, err
	}

	internal.SetContentHeaders(&req.Header)

	res, err := cl.Do(req)
	if err != nil {
		return Token{}, err
	}

	if res.StatusCode != http.StatusOK {
		return Token{}, internal.HTTPError(res)
	}

	var token Token
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

// Refresh an existing token against the API server and return the whole response
// from the API, with the intention of forwarding it back to the client application
func refreshToken(cl *http.Client, config Config, refresh string, scopes []string) (*http.Response, error) {
	values := map[string]string{
		"client_id":     strconv.FormatInt(int64(config.OAuth.ClientId), 10),
		"client_secret": config.OAuth.TokenSecret,
		"grant_type":    "refresh_token",
		"refresh_token": refresh,
		"scope":         strings.Join(scopes, " "),
	}

	body, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		config.TokenURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	internal.SetContentHeaders(&req.Header)

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, internal.HTTPError(res)
	}

	// Forward response back to the calling function
	return res, nil
}
