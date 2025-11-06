package main

import (
	"encoding/json"
	"os"

	"github.com/Sprinter05/osu-chat/api"
)

// TODO: Add arguments
const CONFIG_FILE = "config/config.json"
const DEFAULT_PERMS = 0755

type Config struct {
	OAuth api.OAuth  `json:"oauth"`
	Token *api.Token `json:"token,omitempty"`
}

func getConfig() (config Config, err error) {
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		return config, err
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&config)
	return config, nil
}

func saveConfig(config Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(CONFIG_FILE, data, DEFAULT_PERMS)
	return err
}
