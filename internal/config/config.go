package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
)

const DefaultBaseURL = "https://open-agent.liepin.com"

type RuntimeConfig struct {
	Token       string
	TokenSource string // "cli", "env", "config"
	BaseURL     string
	Timeout     time.Duration
	Output      string // "pretty" or "json"
}

type MissingTokenError struct{}

func (e *MissingTokenError) Error() string {
	return "missing x-user-token. Provide it via --token, LIEPIN_USER_TOKEN, or CLI config file. Visit https://www.liepin.com/mcp/auth to generate one."
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

func ResolveConfig(tokenFlag, baseURLFlag, outputFlag string, timeoutFlag float64) (*RuntimeConfig, error) {
	output := outputFlag
	if output == "" {
		output = "pretty"
	}
	if output != "pretty" && output != "json" {
		return nil, &InvalidInputError{"--output must be: pretty, json"}
	}

	timeout := timeoutFlag
	if timeout <= 0 {
		timeout = 30.0
	}

	var resolvedToken, source string

	if tokenFlag != "" {
		resolvedToken = tokenFlag
		source = "cli"
	} else {
		envToken := os.Getenv("LIEPIN_USER_TOKEN")
		if envToken != "" {
			resolvedToken = envToken
			source = "env"
		} else {
			stored, err := authstore.LoadToken()
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			if stored == "" {
				return nil, &MissingTokenError{}
			}
			resolvedToken = stored
			source = "config"
		}
	}

	baseURL := baseURLFlag
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &RuntimeConfig{
		Token:       resolvedToken,
		TokenSource: source,
		BaseURL:     baseURL,
		Timeout:     time.Duration(timeout * float64(time.Second)),
		Output:      output,
	}, nil
}

func IsProduction(baseURL string) bool {
	return baseURL == "" || baseURL == DefaultBaseURL
}
