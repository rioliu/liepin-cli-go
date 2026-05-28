package cmd

import (
	"fmt"

	"github.com/rioliu/liepin-cli-go/internal/client"
	"github.com/rioliu/liepin-cli-go/internal/config"
	"github.com/rioliu/liepin-cli-go/internal/onboarding"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "First-time config or refresh login state.",
	Long:  "Interactive auth flow: open browser → paste token → verify and save to local config.",
	Run: func(cmd *cobra.Command, _ []string) {
		verifyToken := func(token string) error {
			cfg, err := config.ResolveConfig(token, config.DefaultBaseURL, "pretty", 30.0)
			if err != nil {
				return err
			}
			c := client.New(client.Config{
				Token:   cfg.Token,
				BaseURL: cfg.BaseURL,
				Timeout: cfg.Timeout,
			})
			_, err = c.Get("/mcp/get-resume")
			return err
		}

		token, err := onboarding.RunAuthSetup(verifyToken)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
			return
		}
		if token != "" {
			fmt.Println("Token saved and verified successfully.")
		}
	},
}
