package output

import (
	"encoding/json"
	"fmt"
	"os"
)

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
