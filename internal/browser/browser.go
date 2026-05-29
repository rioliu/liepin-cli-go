// Package browser launches the system default web browser to the Liepin
// authentication page so users can generate an x-user-token.
package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// AuthURL is the Liepin web page where users can generate an x-user-token.
const AuthURL = "https://www.liepin.com/mcp/auth"

// OpenAuthPage launches the system default browser pointed at AuthURL so the
// user can authenticate. It returns an error if the host OS is not supported
// or if the browser process cannot be started; in those cases the caller
// should instruct the user to visit AuthURL manually.
func OpenAuthPage() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", AuthURL)
	case "linux":
		cmd = exec.Command("xdg-open", AuthURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", AuthURL)
	default:
		return fmt.Errorf("failed to open auth page; please visit %s manually", AuthURL)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open auth page; please visit %s manually", AuthURL)
	}
	return nil
}
