package oauth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sprinter05/osu-chat/internal"
)

const OAUTH_STATE_LENGTH int = 128
const OAUTH_PORT_LENGTH int = 5

type OAuth struct {
	ClientId    int    `json:"client_id"`
	TokenSecret string `json:"token_secret"`
}

/* STATE */

// Creates a random state string with the port bundled in it
// It is given in base64 format
func CreateState() (uint16, string, error) {
	// Create random port excluding well known ports
	port := mrand.IntN((2 ^ 16) - ((2 ^ 10) + 1))
	portStr := fmt.Sprintf("%05d", port)

	// Add random data
	data := make([]byte, 0, OAUTH_STATE_LENGTH)
	max := OAUTH_STATE_LENGTH - OAUTH_PORT_LENGTH
	if _, err := io.ReadAtLeast(rand.Reader, data, max); err != nil {
		return 0, "", err
	}

	// Add port for client callback
	data = append(data, []byte(portStr)...)
	return uint16(port), base64.StdEncoding.EncodeToString(data), nil
}

// Returns the port asocciated to a random state string
func GetPortFromState(state string) (uint16, error) {
	// Decode base64
	base, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return 0, err
	}

	// Get port at the end
	max := OAUTH_STATE_LENGTH - OAUTH_PORT_LENGTH
	portRange := base[max:]
	port, err := strconv.ParseUint(string(portRange), 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(port), nil
}

/* CLIENT CALLBACK */

type TokenResponse struct {
	Token Token  `json:"token"`
	State string `json:"state"`
}

func clientRequest(response TokenResponse) (*http.Request, error) {
	port, err := GetPortFromState(response.State)
	if err != nil {
		return nil, err
	}

	body := make([]byte, 0)
	if response.State == "" {
		body, err = json.Marshal(response.Token)
	} else {
		body, err = json.Marshal(response)
	}

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://localhost:%d/token", port),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	internal.SetContentHeaders(&req.Header)

	return req, nil
}

/* SERVER */

type serverFunc func(w http.ResponseWriter, r *http.Request)

func authorization(client *http.Client, config Config) serverFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		if !url.Has("code") || !url.Has("state") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		code := url.Get("code")
		state := url.Get("state")

		token, err := requestToken(client, config, code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req, err := clientRequest(TokenResponse{
			Token: token,
			State: state,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, req.URL.String(), http.StatusPermanentRedirect)
	}
}

func refreshing(client *http.Client, config Config) serverFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		if !url.Has("refresh") || !url.Has("scopes") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		refresh := url.Get("refresh")
		scopes := strings.Split(url.Get("scopes"), "+")

		token, err := refreshToken(client, config, refresh, scopes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req, err := clientRequest(TokenResponse{
			Token: token,
			State: "",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, req.URL.String(), http.StatusPermanentRedirect)
	}
}

func ServerCallback(config Config) {
	cl := internal.DefaultClient()

	mux := http.NewServeMux()
	mux.HandleFunc("/authorization", authorization(cl, config))
	mux.HandleFunc("/refreshing", refreshing(cl, config))

	srv := &http.Server{
		Addr:    config.Address,
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("OAuth server callback error: %s", err)
	}
}

/* CLIENT */

// Ran by the client application
func ClientCallback(state string, port uint16, send chan<- Token) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		response := new(TokenResponse)
		err := json.NewDecoder(r.Response.Body).Decode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if state != response.State {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		send <- response.Token
		fmt.Fprint(w, "Token generated! You can now close this window!")
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("OAuth client callback error: %s", err)
	}
}
