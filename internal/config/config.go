package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Url             string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfg := Config{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("can't find user home directory: %w", err)
	}

	path := homeDir + "/gatorconfig.json"
	jsonFile, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("can't open gatorconfig file: %w", err)
	}

	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling json data to struct: %w", err)
	}

	return cfg, nil
}
