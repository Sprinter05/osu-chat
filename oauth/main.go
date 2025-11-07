package oauth

import (
	"flag"
	"log"

	"github.com/Sprinter05/osu-chat/internal"
)

var configFile string

type Config struct {
	OAuth       OAuth    `json:"oauth"`
	TokenURL    string   `json:"token_url"`
	CallbackURL string   `json:"callback_url"`
	Address     string   `json:"address"`
	Scopes      []string `json:"scopes"`
}

// Ran at startup
func init() {
	flag.StringVar(
		&configFile, "config", "oauth.json",
		"Configuration file to load, must be in JSON format.",
	)
	flag.Parse()
}

func main() {
	conf := new(Config)
	err := internal.GetConfig(configFile, conf)
	if err != nil {
		log.Fatal(err)
	}

	// Run
	ServerCallback(*conf)
}
