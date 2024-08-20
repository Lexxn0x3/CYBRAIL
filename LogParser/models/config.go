package models

import (
    "encoding/json"
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
    file, err := os.Create(path)
    if err != nil {
 	   return err
    }
    defer file.Close()

    err = json.NewEncoder(file).Encode(config)
    if err != nil {
 	   return err
    }
    return nil
}

