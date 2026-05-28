# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test

```bash
make build              # go build -o bin/liepin-cli
make test               # go test ./... (unit tests, 103 Ginkgo specs)
make test-e2e           # e2e tests (requires LIEPIN_USER_TOKEN)
make test-e2e-verbose   # e2e tests with verbose output
make clean              # remove bin/
```

E2E tests run against the live API. Read-only endpoints (`get-resume`, `search-job`) always execute. Write endpoints check `config.IsProduction(baseURL)` and skip when targeting the default production URL. Set `LIEPIN_BASE_URL` to a test environment to run the full suite.

## Architecture

### Command layer (`cmd/`)

Cobra CLI with package-level flag variables registered in `init()`. Each subcommand's `Run` follows the same flow:

1. Build a `map[string]any` of overrides from parsed flags
2. Call `buildPayload(overrides, model.Validate)` — loads optional `--input` JSON file, merges with overrides, validates
3. Call `executeGet` or `executePost` — creates client via `buildClient()`, sends request, renders output

`cmd/common.go` holds shared helpers: `buildClient`, `executeGet`, `executePost`, `buildPayload`, `handleError`, `addCommonFlags`.

### Client (`internal/client/`)

Thin HTTP wrapper. `Get`/`Post` attach the `x-user-token` header and parse JSON responses. Non-2xx responses return typed errors (`AuthorizationError` for 401/403, `RequestError` for 4xx+).

### Config (`internal/config/`)

`ResolveConfig` resolves token from: `--token` flag > `LIEPIN_USER_TOKEN` env var > `~/.config/liepin-cli/config.json` (written by `setup`). `DefaultBaseURL` is `https://open-agent.liepin.com`.

### Input models (`internal/models/`)

Each API input is a struct with a `Validate() error` method. `ParseOptionalInt` returns `*int` from a string — errors must be handled, never discarded with `, _ :=`. `StrPtr` returns `*string` for optional string fields.

### Error handling

`handleError` type-switches on custom error types: `client.AuthorizationError` → "refresh your token", `config.MissingTokenError` → "run setup", `client.RequestError` → exit code 1, others → exit code 2.

## Test Conventions

- **Unit tests**: `*_test.go` alongside source using Ginkgo v2 + Gomega. Mock HTTP server in `cmd/cmd_test.go` tests full CLI flows via `rootCmd.SetArgs()` + `rootCmd.Execute()`.
- **E2E tests**: `test/e2e/e2e_test.go` with build tag `e2e`. Uses `client.New` directly (not CLI), asserts response shape.
- Flag variables leak between tests in the same package — `cmd_test.go` has a `resetAllFlags()` helper called in `BeforeEach`.
