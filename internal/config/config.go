package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Url             string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(cfg)
}

func Read() (Config, error) {
	cfg := Config{}

	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("can't get config file path: %w", err)
	}

	jsonFile, err := os.Open(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("can't open gatorconfig file: %w", err)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding data: %w", err)
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("can't find user home directory: %w", err)
	}

	path := filepath.Join(homeDir, configFileName)
	return path, nil
}

func write(cfg *Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("can't get config file path on write: %w", err)
	}

	jsonFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	if err = encoder.Encode(cfg); err != nil {
		return fmt.Errorf("can't write encoded data to file: %w", err)
	}

	return nil
}
