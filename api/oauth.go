package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/browser"
)

const OSU_URL_OAUTH string = "https://osu.ppy.sh/oauth"
const CALLBACK_OAUTH string = "localhost:43200"
const SCOPES_OAUTH string = "chat.read chat.write chat.write_manage friends.read identify"

type OAuth struct {
	ClientId    int    `json:"client_id"`
	TokenSecret string `json:"token_secret"`
}

func oauthState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func oauthServer(state string, c chan<- string) {
	srv := &http.Server{
		Addr: CALLBACK_OAUTH,
	}

	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		if !url.Has("code") || !url.Has("state") {
			fmt.Fprint(w, "Invalid request! You should not be here!")
			return
		}

		code := url.Get("code")
		check := url.Get("state")

		if state != check {
			fmt.Fprint(w, "Invalid state given!")
			return
		}

		c <- code
		fmt.Fprint(w, "OAuth code obtained! Go back to the application")
		srv.Close()
	})

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("OAuth server error: %s", err)
	}
}

func oauthAuthorize(oauth OAuth) (string, error) {
	state, err := oauthState(128)
	if err != nil {
		return "", err
	}

	url := url.Values{}
	url.Add("client_id", strconv.FormatInt(int64(oauth.ClientId), 10))
	url.Add("state", state)
	url.Add("response_type", "code")
	url.Add("redirect_uri", fmt.Sprintf("http://%s/oauth", CALLBACK_OAUTH))
	url.Add("scope", SCOPES_OAUTH)
	params := url.Encode()

	full := fmt.Sprintf("https://%s/authorize?%s", OSU_URL_OAUTH, params)

	channel := make(chan string)
	go oauthServer(state, channel)

	browser.OpenURL(full)
	code := <-channel

	return code, nil
}
