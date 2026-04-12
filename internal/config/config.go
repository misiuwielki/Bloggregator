package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Url              string `json:"db_url"`
	Current_username string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getCfgPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("Error while getting home directory: %w", err)
		return "", err
	}
	path := filepath.Join(home, configFileName)
	return path, nil
}

func write(cfg Config) error {
	jsn, err := json.Marshal(cfg)
	if err != nil {
		err = fmt.Errorf("Error while marshaling config file: %w", err)
		return err
	}
	path, err := getCfgPath()
	if err != nil {
		err = fmt.Errorf("While writing: %w", err)
		return err
	}
	err = os.WriteFile(path, jsn, 0644)
	if err != nil {
		err = fmt.Errorf("Error while overwriting config file: %w", err)
		return err
	}
	return nil

}

func Read() (Config, error) {
	path, err := getCfgPath()
	if err != nil {
		err = fmt.Errorf("While reading: %w", err)
		return Config{}, err
	}
	jsn, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error while reading config file: %w", err)
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(jsn, &cfg)
	if err != nil {
		err = fmt.Errorf("Error while unmarshaling config file: %w", err)
		return Config{}, err
	}
	return cfg, nil
}

func (cfg Config) SetUser(username string) error {
	cfg.Current_username = username
	err := write(cfg)
	if err != nil {
		err = fmt.Errorf("While setting user: %w", err)
	}
	return nil
}
