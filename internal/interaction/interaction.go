package interaction

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
