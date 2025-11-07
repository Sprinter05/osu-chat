package main

import (
	"flag"
	"log"

	"github.com/Sprinter05/osu-chat/internal"
	"github.com/Sprinter05/osu-chat/internal/conf"
	"github.com/Sprinter05/osu-chat/ui"
)

/* SETUP */

var configFile string

func init() {
	flag.StringVar(
		&configFile, "config", "config.json",
		"Configuration file to load, must be in JSON format.",
	)
	flag.Parse()
}

func main() {
	config := new(conf.Config)
	client := internal.DefaultClient()
	err := conf.GetConfig(configFile, config)
	if err != nil {
		log.Fatal(err)
	}

	gui := ui.GUI{
		Client: client,
	}

	gui.Run()
}
