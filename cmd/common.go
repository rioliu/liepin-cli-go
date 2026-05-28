package cmd

import (
	"fmt"
	"os"

	"github.com/rioliu/liepin-cli-go/internal/client"
	"github.com/rioliu/liepin-cli-go/internal/config"
	"github.com/rioliu/liepin-cli-go/internal/models"
	"github.com/rioliu/liepin-cli-go/internal/output"
	"github.com/rioliu/liepin-cli-go/internal/payload"
	"github.com/spf13/cobra"
)

var (
	tokenFlag    string
	baseURLFlag  string
	timeoutFlag  float64
	outputFlag   string
	inputFlag    string
)

func buildClient() (*client.Client, *config.RuntimeConfig, error) {
	cfg, err := config.ResolveConfig(tokenFlag, baseURLFlag, outputFlag, timeoutFlag)
	if err != nil {
		return nil, nil, err
	}
	return client.New(client.Config{
		Token:   cfg.Token,
		BaseURL: cfg.BaseURL,
		Timeout: cfg.Timeout,
	}), cfg, nil
}

func exitWithError(err error, isRequestError bool) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	code := 2
	if isRequestError {
		code = 1
	}
	os.Exit(code)
}

func handleError(err error) {
	if err == nil {
		return
	}
	switch e := err.(type) {
	case *client.RequestError:
		exitWithError(e, true)
	case *client.AuthorizationError:
		fmt.Fprintln(os.Stderr, "Authorization failed. Please run liepin-cli setup or liepin-cli auth setup to refresh your token.")
		exitWithError(e, true)
	case *config.MissingTokenError:
		fmt.Fprintln(os.Stderr, "No token found. Please run the setup flow first.")
		exitWithError(e, false)
	case *config.InvalidInputError:
		exitWithError(e, false)
	default:
		exitWithError(e, false)
	}
}

func executeGet(path string) {
	c, cfg, err := buildClient()
	if err != nil {
		handleError(err)
		return
	}

	data, err := c.Get(path)
	if err != nil {
		handleError(err)
		return
	}

	output.Render(data, cfg.Output)
}

func executePost(path string, payloadData map[string]any) {
	c, cfg, err := buildClient()
	if err != nil {
		handleError(err)
		return
	}

	data, err := c.Post(path, payloadData)
	if err != nil {
		handleError(err)
		return
	}

	output.Render(data, cfg.Output)
}

func buildPayload(overrides map[string]any, validate func() error) (map[string]any, error) {
	base, err := payload.LoadPayloadFile(inputFlag)
	if err != nil {
		return nil, err
	}

	merged := payload.MergePayload(base, overrides)

	if err := validate(); err != nil {
		return nil, err
	}

	return merged, nil
}

func setIfNotEmpty(m map[string]any, key, value string) {
	if value != "" {
		m[key] = value
	}
}

func setIntIfNotEmpty(m map[string]any, key, raw string) error {
	if raw == "" {
		return nil
	}
	v, err := models.ParseOptionalInt(raw)
	if err != nil {
		return err
	}
	if v != nil {
		m[key] = *v
	}
	return nil
}

func optionalIntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&tokenFlag, "token", "", "Liepin user token")
	cmd.Flags().StringVar(&baseURLFlag, "base-url", "", "API base URL")
	cmd.Flags().Float64Var(&timeoutFlag, "timeout", 0, "Request timeout in seconds")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output format (pretty/json)")
	cmd.Flags().StringVar(&inputFlag, "input", "", "JSON request body file path")

	// Hide advanced flags
	cmd.Flags().MarkHidden("token")
	cmd.Flags().MarkHidden("base-url")
	cmd.Flags().MarkHidden("timeout")
	cmd.Flags().MarkHidden("input")
}
