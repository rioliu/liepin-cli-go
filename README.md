# liepin-cli

A CLI tool for [Liepin](https://www.liepin.com) resume management and job applications.

## Prerequisites

- Go 1.26+
- A Liepin user token from [https://www.liepin.com/mcp/auth](https://www.liepin.com/mcp/auth)

## Quick Start

```bash
git clone https://github.com/rioliu/liepin-cli-go.git
cd liepin-cli-go
make build
export PATH="$PWD/bin:$PATH"
liepin-cli setup
```

Add the `export PATH` line to `~/.zshrc` or `~/.bashrc` to make it permanent.

`setup` opens your browser to the Liepin auth page. Paste your token when prompted — it's verified against the API and saved to `~/.config/liepin-cli/config.json`.

## Token Sources

The CLI resolves your token in this order (first found wins):

1. `--token` flag
2. `LIEPIN_USER_TOKEN` environment variable
3. `~/.config/liepin-cli/config.json` (written by `setup`)

## Commands

### Authentication

```bash
liepin-cli auth status     # Show current login status
liepin-cli auth clear      # Remove saved token
liepin-cli auth open       # Open auth page in browser
liepin-cli auth setup      # Interactive setup (same as liepin-cli setup)
```

### Resume

```bash
# Read
liepin-cli resume get

# Base info
liepin-cli resume update-base-info --real-name "Zhang San" --sex 1 --city-code 010

# Self assessment
liepin-cli resume update-self-assess --self-assess "Experienced backend engineer..."

# Education
liepin-cli resume add-edu-exp --school "Peking University" --start 201909 --end 202306 --degree 5
liepin-cli resume update-edu-exp --edu-id 12345 --school "Peking University"

# Work experience
liepin-cli resume add-work-exp --comp-name "Acme Corp" --work-start 202001 --work-end 202406
liepin-cli resume update-work-exp --work-id 12345 --comp-name "Acme Corp"

# Project experience
liepin-cli resume add-project-exp --name "Migration Platform" --start 202301 --end 202312
liepin-cli resume update-project-exp --id 12345 --name "Migration Platform v2"

# Job preferences
liepin-cli resume add-job-want --jobtitle "Backend Engineer" --dq "北京"
liepin-cli resume update-job-want --id 12345 --jobtitle "Senior Backend Engineer"
```

### Job

```bash
liepin-cli job search --job-name "后端开发" --address "北京" --page 0
liepin-cli job apply --job-id 123456 --job-kind 1
```

### Global Flags

| Flag | Description |
|------|-------------|
| `--output`, `-o` | Output format: `pretty` (default) or `json` |
| `--token` | Liepin user token |
| `--base-url` | API base URL (default: `https://open-agent.liepin.com`) |
| `--timeout` | Request timeout in seconds (default: `30`) |
| `--input` | JSON file for request body (merged with flag values) |

## Build & Test

```bash
make build           # Build binary → bin/liepin-cli
make install         # Build and copy to ~/.local/bin (optional, or just add ./bin to PATH)
make test            # Run all unit tests
make test-verbose    # Run unit tests with verbose output
make test-e2e        # Run e2e tests against the real API
make test-e2e-verbose # Run e2e tests with verbose output
make clean           # Remove build artifacts
```

### E2E Tests

E2E tests run against the live Liepin API and require:

```bash
export LIEPIN_USER_TOKEN="your-token"
make test-e2e
```

Read-only endpoints (`get-resume`, `search-job`) always run. Write endpoint tests (apply, update, add) are **skipped when targeting the production API** to avoid creating dirty data. To run the full suite, point at a test environment:

```bash
export LIEPIN_BASE_URL="https://test-env.example.com"
export LIEPIN_USER_TOKEN="your-token"
make test-e2e
```

## Claude Code Integration

To teach Claude how to use `liepin-cli`, install the official skill from the upstream repo:

```bash
npx skills add liepin-tech-2026/liepin-cli -y -g
```

If you don't have Node.js, clone the repo and use `uv`:

```bash
git clone https://github.com/liepin-tech-2026/liepin-cli.git /tmp/liepin-cli
cd /tmp/liepin-cli
uv run liepin-cli skill install --claude
```

