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

// Authentication

func setGenericHeaders(hd *http.Header, token Token) {
	hd.Set("Content-Type", "application/json")
	hd.Set("Accept", "application/json")
	hd.Set("Authorization", fmt.Sprintf(
		"%s %s",
		token.TokenType,
		token.AccessToken,
	))
}

func GetChannelList(cl *http.Client, token Token) ([]ChatChannel, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		OsuApiUrl+"/chat/channels",
		nil,
	)

	if err != nil {
		return nil, err
	}

	setGenericHeaders(&req.Header, token)

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, internal.HTTPError(res)
	}

	list := make([]ChatChannel, 0)
	err = json.NewDecoder(res.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
