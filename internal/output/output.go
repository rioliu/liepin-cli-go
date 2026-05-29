// Package output renders command results to stdout in either a human
// friendly "pretty" form or a machine-readable JSON form.
package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Render prints data to stdout using the given mode ("pretty" or "json").
// In pretty mode, a nil value prints a generic success message, strings are
// printed verbatim, and structured values are still rendered as indented
// JSON for readability.
func Render(data any, mode string) {
	if mode == "json" {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
		return
	}

	// pretty mode
	if data == nil {
		fmt.Println("Success.")
		return
	}

	switch v := data.(type) {
	case string:
		fmt.Println(v)
	default:
		jsonData, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	}
}
