package api

import (
	"encoding/json"
	"net/http"
)

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

	var list []ChatChannel
	err = json.NewDecoder(res.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
