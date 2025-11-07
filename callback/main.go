package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sprinter05/osu-chat/callback/oauth"
	"github.com/Sprinter05/osu-chat/internal/conf"
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
	log.SetOutput(os.Stdout)
	config := new(oauth.Config)
	err := conf.GetConfig(configFile, config)
	if err != nil {
		log.Fatal(err)
	}

	// Run
	oauth.ServerCallback(*config)

	fmt.Print("Server terminated!")
}
