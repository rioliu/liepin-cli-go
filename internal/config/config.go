// Package config resolves runtime configuration for the CLI by merging
// command-line flags, environment variables, and on-disk settings.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
)

// DefaultBaseURL is the production Liepin open-agent API endpoint used when
// the user does not supply --base-url.
const DefaultBaseURL = "https://open-agent.liepin.com"

// RuntimeConfig is the fully resolved configuration passed to commands. It
// captures the chosen authentication token along with its origin so that
// commands can produce informative diagnostics.
type RuntimeConfig struct {
	Token       string
	TokenSource string // "cli", "env", "config"
	BaseURL     string
	Timeout     time.Duration
	Output      string // "pretty" or "json"
}

// MissingTokenError is returned by ResolveConfig when no x-user-token can be
// found via any of the supported sources.
type MissingTokenError struct{}

func (e *MissingTokenError) Error() string {
	return "missing x-user-token. Provide it via --token, LIEPIN_USER_TOKEN, or CLI config file. Visit https://www.liepin.com/mcp/auth to generate one."
}

// InvalidInputError represents a user-supplied flag value that fails
// validation (for example an unsupported --output mode).
type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

// ResolveConfig builds a RuntimeConfig from CLI flags, falling back to the
// LIEPIN_USER_TOKEN environment variable and the persisted config file when
// --token is not set. It returns *InvalidInputError for malformed flags and
// *MissingTokenError when no token can be located.
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

// IsProduction reports whether the given base URL refers to the production
// Liepin endpoint (or has been left at the default).
func IsProduction(baseURL string) bool {
	return baseURL == "" || baseURL == DefaultBaseURL
}
