package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Students []Student `json:"students"`
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func SaveConfig(path string, config Config) error {
	// Create or overwrite the config file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	// Create a JSON encoder with indentation for readability
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// Encode the config to JSON and write to the file
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
