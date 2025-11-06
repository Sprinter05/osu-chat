package api

import (
	"encoding/json"
	"net/http"
)

func GetChannelList(cl *http.Client, oauth OAuth) ([]ChatChannel, error) {
	req, err := http.NewRequest("GET", OSU_URL+"/chat/channels", nil)
	if err != nil {
		return nil, err
	}

	setGenericHeaders(&req.Header, oauth)

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
