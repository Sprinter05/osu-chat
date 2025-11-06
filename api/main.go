package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const OSU_URL string = "https://osu.ppy.sh/api/v2"

type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RequestToken(cl *http.Client, oauth OAuth) (Token, error) {
	code, err := oauthAuthorize(oauth)
	if err != nil {
		return Token{}, err
	}

	url := url.Values{}
	url.Add("client_id", strconv.FormatInt(int64(oauth.ClientId), 10))
	url.Add("client_secret", oauth.TokenSecret)
	url.Add("code", code)
	url.Add("grant_type", "authorization_code")

	res, err := cl.PostForm(OSU_URL_OAUTH+"/token", url)
	if err != nil {
		return Token{}, err
	}

	var token Token
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func DeleteToken(cl *http.Client, token Token) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		OSU_URL_OAUTH+"/tokens/current",
		nil,
	)

	setGenericHeaders(&req.Header, token)

	_, err = cl.Do(req)
	return err
}
