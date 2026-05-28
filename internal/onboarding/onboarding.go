package onboarding

import (
	"fmt"
	"os"
	"strings"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
	"github.com/rioliu/liepin-cli-go/internal/browser"
	"github.com/rioliu/liepin-cli-go/internal/interaction"
)

func MaskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return "****"
	}
	if idx := strings.LastIndex(token, "-"); idx >= 0 {
		head := token[:idx]
		tail := token[idx+1:]
		if len(tail) <= 4 {
			return head + "-****"
		}
		return head + "-****" + tail[len(tail)-4:]
	}
	return "****" + token[len(token)-4:]
}

func RunAuthSetup(verify func(string) error) (string, error) {
	interactive := interaction.IsInteractiveTerminal()

	if interactive && interaction.ConfirmOpenAuthPage() {
		if err := browser.OpenAuthPage(); err != nil {
			fmt.Fprintf(os.Stderr, "Note: %v\n", err)
		}
	}

	token := interaction.PromptForToken()
	if token == "" {
		if interactive {
			return "", fmt.Errorf("no token received (empty or cancelled). Please paste the token, or run liepin-cli auth open first, then re-run setup")
		}
		return "", fmt.Errorf("setup requires an interactive terminal that supports paste input")
	}

	if err := verify(token); err != nil {
		return "", err
	}

	if err := authstore.SaveToken(token); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return token, nil
}
