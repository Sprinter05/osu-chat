package internal

import (
	"encoding/json"
	"os"
)

const DefaultPerms = 0755

// Gets the config into the "config" variable.
// Said variable must be a pointer
func GetConfig(path string, config any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(config)
	return nil
}

// Saves the config into the specified file.
// The provided configuration must not be a pointer
func SaveConfig(path string, config any) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, DefaultPerms)
	return err
}
