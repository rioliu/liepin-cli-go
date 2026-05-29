// Package authstore persists the user's x-user-token to a JSON config file
// under $XDG_CONFIG_HOME/liepin-cli (or ~/.config/liepin-cli) so subsequent
// invocations of the CLI can reuse it without prompting.
package authstore

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type storedConfig struct {
	Token string `json:"token"`
}

func configPath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "liepin-cli", "config.json"), nil
}

func readConfig() (storedConfig, error) {
	path, err := configPath()
	if err != nil {
		return storedConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return storedConfig{}, nil
		}
		return storedConfig{}, err
	}

	var cfg storedConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return storedConfig{}, err
	}
	return cfg, nil
}

func writeConfig(cfg storedConfig) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadToken reads the stored authentication token from the local config file.
// It returns an empty string if no token has been saved yet.
func LoadToken() (string, error) {
	cfg, err := readConfig()
	if err != nil {
		return "", err
	}
	return cfg.Token, nil
}

// SaveToken persists the given authentication token to the local config file,
// creating or updating the file as needed.
func SaveToken(token string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	cfg.Token = token
	return writeConfig(cfg)
}

// ClearToken removes the stored authentication token from the config file.
// It reports whether a token was actually cleared (true) or if there was
// nothing to remove (false).
func ClearToken() (bool, error) {
	path, err := configPath()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}

	cfg, err := readConfig()
	if err != nil {
		return false, err
	}

	if cfg.Token == "" {
		return false, nil
	}

	cfg.Token = ""
	return true, writeConfig(cfg)
}
