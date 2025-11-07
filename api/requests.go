package api

import (
	"encoding/json"
	"net/http"
)

type Token struct {
	TokenType    string
	AccessToken  string
	RefreshToken string
}

func GetChannelList(cl *http.Client, token Token) ([]ChatChannel, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		OSU_URL+"/chat/channels",
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
		return nil, httpErr(res)
	}

	list := make([]ChatChannel, 0)
	err = json.NewDecoder(res.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
