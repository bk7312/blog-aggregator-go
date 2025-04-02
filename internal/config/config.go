package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	path, err := getConfigFilePath()
	if err != nil {
		log.Fatal("config file path not found")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("unable to open config file")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("error reading config")
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("error unmarshalling config")
	}

	return cfg
}

func SetUser(name string) {
	cfg := Read()
	cfg.CurrentUserName = name
	err := write(cfg)
	if err != nil {
		log.Fatal("error setting username")
	}
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, configFileName)
	return path, nil
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal("error marshal config")
	}
	path, err := getConfigFilePath()
	if err != nil {
		log.Fatal("error getting config path")
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal("error writing to file")
	}
	return nil
}

const configFileName = ".gatorconfig.json"
