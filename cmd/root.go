// Package cmd assembles the cobra command tree for the liepin-cli binary
// and exposes Execute as the single entry point used by main.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "liepin-cli",
	Short: "Liepin resume and job CLI.",
	Long:  "Liepin resume and job CLI — local CLI for Liepin resume queries/updates, job search and applications.",
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(jobCmd)
}

// Execute runs the root cobra command and returns any error produced by the
// selected sub-command. It is the single entry point used by main.
func Execute() error {
	return rootCmd.Execute()
}
