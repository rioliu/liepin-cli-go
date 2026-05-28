package output

import (
	"encoding/json"
	"fmt"
)

func Render(data any, mode string) {
	if mode == "json" {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
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
		jsonData, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(jsonData))
	}
}
