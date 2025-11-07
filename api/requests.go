package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sprinter05/osu-chat/internal"
)

type OAuth struct {
	ClientId    int      `json:"client_id"`
	CallbackURL string   `json:"callback_url"`
	Scopes      []string `json:"scopes"`
}

type Token struct {
	TokenType    string
	AccessToken  string
	RefreshToken string
}

/* AUXILIARY FUNCTIONS */

// Sets the generic headers needed for an API request
func setGenericHeaders(hd *http.Header, token Token) {
	internal.SetContentHeaders(hd)
	hd.Set("Authorization", fmt.Sprintf(
		"%s %s",
		token.TokenType,
		token.AccessToken,
	))
}

// Base function to do a request to the osu API without a body
//
// method: HTTP method,
// endpoint: URL route after the base address (e.g "/a/b"),
// expected: HTTP code that should be returned by the request
func makeRequestNoBody(cl *http.Client, method string, endpoint string, expected int, token Token) (*http.Response, error) {
	req, err := http.NewRequest(
		method, OsuApiUrl+endpoint, nil,
	)
	if err != nil {
		return nil, err
	}

	setGenericHeaders(&req.Header, token)

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != expected {
		return nil, internal.HTTPError(res)
	}

	return res, nil
}

/* API REQUESTS */

// Gets the list of all public chat channels
func GetChannelList(cl *http.Client, token Token) ([]ChatChannel, error) {
	res, err := makeRequestNoBody(
		cl, http.MethodGet,
		"/chat/channels",
		http.StatusOK, token,
	)
	if err != nil {
		return nil, err
	}

	list := make([]ChatChannel, 0)
	err = json.NewDecoder(res.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Removes an existing token that is currently available
func DeleteToken(cl *http.Client, token Token) error {
	_, err := makeRequestNoBody(
		cl, http.MethodDelete,
		"/oauth/tokens/current",
		http.StatusNoContent, token,
	)

	return err
}
