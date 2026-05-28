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

func Execute() error {
	return rootCmd.Execute()
}
