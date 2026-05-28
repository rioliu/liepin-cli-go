package interaction

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func IsInteractiveTerminal() bool {
	stdinInfo, _ := os.Stdin.Stat()
	stdoutInfo, _ := os.Stdout.Stat()
	return (stdinInfo.Mode()&os.ModeCharDevice) != 0 &&
		(stdoutInfo.Mode()&os.ModeCharDevice) != 0
}

func ConfirmOpenAuthPage() bool {
	if !IsInteractiveTerminal() {
		return false
	}
	fmt.Print("Open the Liepin auth page now? [Y/n] ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes"
}

func PromptForToken() string {
	if !IsInteractiveTerminal() {
		return ""
	}
	fmt.Print("Please paste x-user-token: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
