package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Sprinter05/osu-chat/callback/oauth"
	"github.com/Sprinter05/osu-chat/internal"
)

var configFile string

func init() {
	flag.StringVar(
		&configFile, "config", "oauth.json",
		"Configuration file to load, must be in JSON format.",
	)
	flag.Parse()
}

func main() {
	conf := new(oauth.Config)
	err := internal.GetConfig(configFile, conf)
	if err != nil {
		log.Fatal(err)
	}

	// Run
	oauth.ServerCallback(*conf)

	fmt.Print("Server terminated!")
}
