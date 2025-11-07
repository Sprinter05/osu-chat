package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sprinter05/osu-chat/internal"
	"github.com/Sprinter05/osu-chat/oauth"
	"github.com/pkg/browser"
)

const OsuApiUrl string = "https://osu.ppy.sh/api/v2"
const OsuOauthUrl string = "https://osu.ppy.sh/oauth/authorize"

func RetrieveToken(params OAuth) (oauth.Token, error) {
	port, state, err := oauth.CreateState()
	if err != nil {
		return oauth.Token{}, err
	}

	ret, _ := oauth.GetPortFromState(state)
	if ret != port {
		panic("invalid")
	}

	url := url.Values{}
	url.Add("client_id", strconv.FormatInt(int64(params.ClientId), 10))
	url.Add("state", state)
	url.Add("response_type", "code")
	url.Add("redirect_uri", params.CallbackURL)
	url.Add("scope", strings.Join(params.Scopes, " "))
	full := url.Encode()

	osuUrl := fmt.Sprintf("%s?%s", OsuOauthUrl, full)

	channel := make(chan oauth.Token)
	go oauth.ClientCallback(state, port, channel)

	browser.OpenURL(osuUrl)
	token := <-channel

	return token, nil
}

func DeleteToken(cl *http.Client, token Token) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		OsuApiUrl+"/oauth/tokens/current",
		nil,
	)

	if err != nil {
		return err
	}

	setGenericHeaders(&req.Header, token)

	res, err := cl.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return internal.HTTPError(res)
	}

	return nil
}
