package cmd

import (
	"fmt"

	"github.com/rioliu/liepin-cli-go/internal/authstore"
	"github.com/rioliu/liepin-cli-go/internal/browser"
	"github.com/rioliu/liepin-cli-go/internal/onboarding"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication.",
	Long:  "View status, clear token, open auth page, interactive setup.",
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current login status.",
	Run: func(cmd *cobra.Command, _ []string) {
		token, err := authstore.LoadToken()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			return
		}
		if token != "" {
			fmt.Printf("Saved token: %s\n", onboarding.MaskToken(token))
		} else {
			fmt.Println("No saved token in config file.")
		}
	},
}

var authClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear current auth info.",
	Run: func(cmd *cobra.Command, _ []string) {
		cleared, err := authstore.ClearToken()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			return
		}
		if cleared {
			fmt.Println("Saved token cleared.")
		} else {
			fmt.Println("No saved token to clear.")
		}
	},
}

var authOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the auth page.",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := browser.OpenAuthPage(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			return
		}
		fmt.Println("Opened Liepin auth page.")
	},
}

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up or update auth.",
	Long:  "Interactive auth flow, equivalent to liepin-cli setup.",
	Run:   setupCmd.Run,
}

func init() {
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authClearCmd)
	authCmd.AddCommand(authOpenCmd)
	authCmd.AddCommand(authSetupCmd)
}
