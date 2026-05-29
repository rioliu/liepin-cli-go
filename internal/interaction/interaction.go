// Package interaction provides small helpers for terminal interactions such
// as detecting an interactive TTY and prompting the user for confirmation
// or sensitive input.
package interaction

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// IsInteractiveTerminal reports whether both stdin and stdout are connected
// to a character device (i.e. a real terminal), so that interactive prompts
// are safe to issue.
func IsInteractiveTerminal() bool {
	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	stdoutInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stdinInfo.Mode()&os.ModeCharDevice) != 0 &&
		(stdoutInfo.Mode()&os.ModeCharDevice) != 0
}

// ConfirmOpenAuthPage prompts the user (Y/n, defaulting to yes) for
// permission to launch the Liepin auth page. It returns false when stdin or
// stdout is not a terminal, or when the user explicitly declines.
func ConfirmOpenAuthPage() bool {
	if !IsInteractiveTerminal() {
		return false
	}
	fmt.Print("Open the Liepin auth page now? [Y/n] ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes"
}

// PromptForToken asks the user to paste an x-user-token into the terminal
// and returns the trimmed value. It returns an empty string when not running
// in an interactive terminal or when reading from stdin fails.
func PromptForToken() string {
	if !IsInteractiveTerminal() {
		return ""
	}
	fmt.Print("Please paste x-user-token: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}
