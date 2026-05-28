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

func LoadToken() (string, error) {
	cfg, err := readConfig()
	if err != nil {
		return "", err
	}
	return cfg.Token, nil
}

func SaveToken(token string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	cfg.Token = token
	return writeConfig(cfg)
}

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
