package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const fileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string
	CurrentUserName string
}

func Read() (Config, error) {
	configFilePath, err := getPathFileConfig()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func SetUsername(username string) error {
	configFilePath, err := getPathFileConfig()
	if err != nil {
		return err
	}

	hasConfigFile := fileExists(configFilePath)
	var config Config

	if hasConfigFile {
		if config, err = Read(); err != nil {
			return err
		}
	} else {
		config = Config{DbUrl: "https://exempleurl.com"}
	}

	config.CurrentUserName = username

	file, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func getPathFileConfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, fileName), nil
}
