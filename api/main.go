package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const OSU_URL string = "https://osu.ppy.sh/api/v2"

type APIToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RequestToken(cl *http.Client, oauth OAuth) (APIToken, error) {
	code, err := oauthAuthorize(oauth)
	if err != nil {
		return APIToken{}, err
	}

	values := map[string]string{
		"client_id":     strconv.FormatInt(int64(oauth.ClientId), 10),
		"client_secret": oauth.TokenSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  fmt.Sprintf("http://%s/oauth", CALLBACK_OAUTH),
	}
	body, err := json.Marshal(values)
	if err != nil {
		return APIToken{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		OSU_URL_OAUTH+"/token",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := cl.Do(req)
	if err != nil {
		return APIToken{}, err
	}

	if res.StatusCode != http.StatusOK {
		return APIToken{}, httpErr(res)
	}

	var token APIToken
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return APIToken{}, err
	}

	return token, nil
}

func DeleteToken(cl *http.Client, token Token) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		OSU_URL+"/oauth/tokens/current",
		nil,
	)

	setGenericHeaders(&req.Header, token)

	res, err := cl.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return httpErr(res)
	}

	return nil
}
