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

const OauthStateLength int = 128
const OauthPortLength int = 5

const MinPort int = 2<<9 + 1
const MaxPort int = 2<<15 - 1

type OAuth struct {
	ClientId    int    `json:"client_id"`
	TokenSecret string `json:"token_secret"`
}

type Config struct {
	OAuth       OAuth  `json:"oauth"`
	TokenURL    string `json:"token_url"`
	CallbackURL string `json:"callback_url"`
	Address     string `json:"address"`
}

/* STATE */

// Creates a random state string with the port at the end of
// the string, given in base64 format
func CreateState() (uint16, string, error) {
	// Create random port excluding well known ports
	port := mrand.IntN(MaxPort-MinPort) + MinPort
	portStr := fmt.Sprintf("%05d", port)

	// Add random data
	random := OauthStateLength - OauthPortLength
	data := make([]byte, random)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return 0, "", err
	}

	// Add port for client callback
	data = append(data, []byte(portStr)...)
	return uint16(port), base64.StdEncoding.EncodeToString(data), nil
}

// Returns the port asocciated to a random state string that
// must be in base64 format
func GetPortFromState(state string) (uint16, error) {
	// Decode base64
	base, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return 0, err
	}

	// Get port at the end of the string
	nonRand := OauthStateLength - OauthPortLength
	portRange := base[nonRand:]
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

// Redirects the token back to the native application, running an http server
// on an ephemeral port. Only used for authorization
func clientRequest(response TokenResponse) (*http.Request, error) {
	port, err := GetPortFromState(response.State)
	if err != nil {
		return nil, err
	}

	// We encode the token in JSON
	body := make([]byte, 0)
	body, err = json.Marshal(response)
	if err != nil {
		return nil, err
	}

	jsonVal := base64.StdEncoding.EncodeToString(body)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://localhost:%d/token?json=%s", port, jsonVal),
		nil,
	)
	if err != nil {
		return nil, err
	}

	internal.SetContentHeaders(&req.Header)

	return req, nil
}

/* SERVER */

type serverFunc func(w http.ResponseWriter, r *http.Request)

// Auxiliary method to log an error on the petition
func logError(endpoint string, r *http.Request, err error) {
	log.Printf(
		"[X] Error on %s endpoint from %s: %s\n",
		endpoint, r.RemoteAddr, err.Error(),
	)
}

// Auxiliary method to log a successful operation.
// It doesn't print the entire state string
func logSuccess(endpoint string, r *http.Request, state string) {
	log.Printf(
		"[i] Successful operation performed on %s endpoint from %s with state %s...\n",
		endpoint, r.RemoteAddr, state[:32],
	)
}

// HTTP method for the "/authorization" endpoint. The API must redirect
// to this endpoint, which will retrieve the token and then send it back
// to the native application through another redirect
func authorization(client *http.Client, config Config) serverFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		if !url.Has("code") || !url.Has("state") {
			logError("authorization", r, fmt.Errorf("no code or state provided"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		code := url.Get("code")
		state := url.Get("state")

		token, err := requestToken(client, config, code)
		if err != nil {
			logError("authorization", r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send token back to the client application
		req, err := clientRequest(TokenResponse{
			Token: token,
			State: state,
		})
		if err != nil {
			logError("authorization", r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, req.URL.String(), http.StatusPermanentRedirect)
		logSuccess("authorization", r, state)
	}
}

// HTTP method for the "/refresh" endpoint. The client app sends a POST request
// giving the refresh token and the scopes, which will then be used to refresh
// the token, returning the API response back to the client app
func refresh(client *http.Client, config Config) serverFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		if !url.Has("refresh") || !url.Has("scopes") {
			logError("refresh", r, fmt.Errorf("no refresh or scopes provided"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		refresh := url.Get("refresh")
		scopes := strings.Split(url.Get("scopes"), " ")

		response, err := refreshToken(client, config, refresh, scopes)
		if err != nil {
			logError("refresh", r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		buf, err := io.ReadAll(response.Body)
		if err != nil {
			logError("refresh", r, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send token back to client
		w.Write(buf)
		logSuccess("refresh", r, refresh)
	}
}

// Function to be ran by the OAuth intermediate server to start
// listening for petitions. The only way to close the HTTP server
// is by closing the running routine
func ServerCallback(config Config) {
	cl := internal.DefaultClient()

	mux := http.NewServeMux()
	mux.HandleFunc("/authorization", authorization(cl, config))
	mux.HandleFunc("/refresh", refresh(cl, config))

	srv := &http.Server{
		Addr:    config.Address,
		Handler: mux,
	}

	fmt.Printf("Started OAuth Intermediate Server on %s\n", config.Address)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("OAuth intermediate server callback error: %s", err)
	}
}

/* CLIENT */

// Function to be ran by the native client app to wait for the OAuth
// intermadiate server to redirect here. It will close once
// the operation has been performed
func ClientCallback(state string, port uint16, send chan<- Token) {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: mux,
	}

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		// TODO: close server according to GUI to show html message
		defer srv.Close()

		url := r.URL.Query()
		if !url.Has("json") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonStr := url.Get("json")
		jsonRaw, err := base64.StdEncoding.DecodeString(jsonStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response TokenResponse
		err = json.NewDecoder(bytes.NewBuffer(jsonRaw)).Decode(&response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Verify againsy cross-forgery
		if state != response.State {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		send <- response.Token
		fmt.Fprint(w, "Token generated! You can now close this window!")
	})

	// Wait until the token has been obtained
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("OAuth client callback error: %s", err)
	}
}
