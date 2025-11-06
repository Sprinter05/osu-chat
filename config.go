package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Sprinter05/osu-chat/api"
)

// TODO: Add arguments
const CONFIG_FILE = "config/config.json"

type Config struct {
	OAuth api.OAuth `json:"oauth"`
}

func getConfig() (config Config) {
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		log.Print("OAuth Token required!")
		log.Fatal(err)
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&config)
	return config
}
