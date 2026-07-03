package cmd

import (
	"fmt"
	"os"

	"github.com/rioliu/liepin-cli-go/internal/client"
	"github.com/rioliu/liepin-cli-go/internal/config"
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
	insecureFlag bool
)

func buildClient() (*client.Client, *config.RuntimeConfig, error) {
	cfg, err := config.ResolveConfig(tokenFlag, baseURLFlag, outputFlag, timeoutFlag)
	if err != nil {
		return nil, nil, err
	}
	return client.New(client.Config{
		Token:              cfg.Token,
		BaseURL:            cfg.BaseURL,
		Timeout:            cfg.Timeout,
		InsecureSkipVerify: insecureFlag,
	}), cfg, nil
}

func exitWithError(err error, isRequestError bool) {
	if outputFlag == "json" {
		errMap := map[string]any{
			"errCode": -1,
			"errMsg":  err.Error(),
		}
		if _, ok := err.(*client.RateLimitError); ok {
			errMap["rateLimited"] = true
		}
		output.Render(errMap, "json")
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
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
		if outputFlag != "json" {
			fmt.Fprintln(os.Stderr, "Authorization failed. Please run liepin-cli setup or liepin-cli auth setup to refresh your token.")
		}
		exitWithError(e, true)
	case *config.MissingTokenError:
		if outputFlag != "json" {
			fmt.Fprintln(os.Stderr, "No token found. Please run the setup flow first.")
		}
		exitWithError(e, false)
	case *config.InvalidInputError:
		exitWithError(e, false)
	case *client.TLSError:
		if outputFlag != "json" {
			fmt.Fprintln(os.Stderr, "TLS verification failed. This can happen when an HTTP proxy intercepts traffic.\n  Rerun with --insecure (-k) to skip certificate verification.")
		}
		exitWithError(e, false)
	case *client.RateLimitError:
		if outputFlag != "json" {
			fmt.Fprintf(os.Stderr, "Rate limited. Please retry after %s.\n", e.RetryAfter)
		}
		exitWithError(e, true)
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
	cmd.Flags().BoolVarP(&insecureFlag, "insecure", "k", false, "Skip TLS certificate verification")

	// Hide advanced flags
	cmd.Flags().MarkHidden("token")
	cmd.Flags().MarkHidden("base-url")
	cmd.Flags().MarkHidden("timeout")
	cmd.Flags().MarkHidden("input")
}
