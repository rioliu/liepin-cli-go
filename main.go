// Command liepin-cli is the entry point for the Liepin command-line tool.
// It delegates to the cmd package which wires up the cobra command tree.
package main

import (
	"fmt"
	"os"

	"github.com/rioliu/liepin-cli-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
