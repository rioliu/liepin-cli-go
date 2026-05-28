# CLAUDE.md

## Build & Test

```bash
make build           # go build -o bin/liepin-cli
make test            # go test ./... (unit tests only)
make test-e2e        # go test -tags=e2e ./test/e2e/ (requires LIEPIN_USER_TOKEN)
make clean           # remove bin/
```

## Test Organization

- **Unit tests**: `*_test.go` alongside source in each package. Mock HTTP server for CLI tests. 103 Ginkgo specs across `cmd/`, `internal/authstore/`, `internal/client/`, `internal/config/`, `internal/models/`, `internal/payload/`.
- **E2E tests**: `test/e2e/e2e_test.go` (build tag `e2e`). Read-only endpoints always run. Write endpoints skip when `config.IsProduction(baseURL)` returns true — pass `LIEPIN_BASE_URL` pointing to a test environment to run the full suite.

## Conventions

- Cobra CLI with package-level flag vars, registered in `init()`.
- Shared helpers live in `cmd/common.go`: `buildClient`, `executeGet`, `executePost`, `buildPayload`, `addCommonFlags`.
- `buildPayload` merges flag overrides into a file-loaded base payload, then validates via model's `Validate()`.
- Token resolution: `--token` flag > `LIEPIN_USER_TOKEN` env > `~/.config/liepin-cli/config.json`.
- Errors from `models.ParseOptionalInt` must be handled — never use `, _ :=` to discard them.
- All user-facing strings in English.
- No Anthropic/AI attribution in git commits or PRs.
