package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

const AuthURL = "https://www.liepin.com/mcp/auth"

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
