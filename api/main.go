package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sprinter05/osu-chat/callback/oauth"
	"github.com/Sprinter05/osu-chat/internal"
	"github.com/pkg/browser"
)

const OsuApiUrl string = "https://osu.ppy.sh/api/v2"
const OsuOauthUrl string = "https://osu.ppy.sh/oauth/authorize"

// Retrieve a token by authorizing the OAuth application and following
// the authorization flow through the intermadiate OAuth server and back
// to the client application
func RetrieveToken(params OAuth) (tok oauth.Token, err error) {
	var port uint16
	var state string

	// Make sure the randomly obtained port is not in use already
	available := false
	for !available {
		port, state, err = oauth.CreateState()
		if err != nil {
			return tok, err
		}

		available = internal.CheckPortInUse(port)
	}

	url := url.Values{}
	url.Add("client_id", strconv.FormatInt(int64(params.ClientId), 10))
	url.Add("state", state)
	url.Add("response_type", "code")
	url.Add("redirect_uri", params.CallbackURL)
	url.Add("scope", strings.Join(params.Scopes, " "))
	full := url.Encode()

	osuUrl := fmt.Sprintf("%s?%s", OsuOauthUrl, full)

	// Run the client loopback server concurrently
	channel := make(chan oauth.Token)
	go oauth.ClientCallback(state, port, channel)

	// Open the authorization in the browser and wait until
	// the token is forwarder back to the client
	browser.OpenURL(osuUrl)
	token := <-channel

	return token, nil
}
