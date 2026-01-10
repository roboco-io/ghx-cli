# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ghx-cli (GitHub eXtensions) is a Go CLI tool for GitHub features not fully supported by the official `gh` CLI. It uses GitHub's GraphQL API to provide:
- **GitHub Projects v2**: View management, workflow automation, bulk operations, analytics
- **GitHub Discussions**: Full discussion management (create, comment, answer, close, lock)

## Build and Development Commands

```bash
# Build
make build              # Build binary to bin/ghx
make install            # Build and install to GOPATH/bin

# Testing
make test               # Run all tests with race detection and coverage
make test-unit          # Run unit tests only (short mode)
go test -v ./internal/service/...  # Run tests for a specific package

# Linting
make lint               # Run golangci-lint (strict config in .golangci.yml)

# Code formatting
make fmt                # Format code with gofmt and tidy modules

# Run directly
go run ./cmd/ghx        # Run without building
```

## Architecture

```
cmd/
  ghx/main.go           # Entry point, sets version info
  root.go               # Root cobra command, registers subcommands

internal/
  api/
    client.go           # GraphQL client with rate limiting and retry logic
    graphql/            # GraphQL queries, mutations, and type definitions
      projects.go       # Project-related queries/mutations
      items.go          # Item operations
      fields.go         # Field operations
      views.go          # View operations
      discussions.go    # Discussion operations
      workflows.go      # Workflow automation
      analytics.go      # Analytics queries

  auth/
    manager.go          # Authentication manager (gh CLI token integration)
    gh_integration.go   # GitHub CLI token extraction

  cmd/                  # Cobra command definitions
    project/            # ghx project [create|list|view|edit|delete|export|import|link|workflow|template]
    item/               # ghx item [add|list|view|edit|remove|add-bulk|update-bulk]
    field/              # ghx field [create|list|update|delete|add-option|update-option|delete-option]
    view/               # ghx view [create|list|update|delete|copy|sort|group]
    discussion/         # ghx discussion [list|view|create|edit|delete|close|reopen|lock|unlock|comment|answer|category]
    analytics/          # ghx analytics [overview|export|bulk-update]
    auth/               # ghx auth [status]

  service/              # Business logic layer between commands and API
    project.go          # Project operations, export/import logic
    item.go             # Item operations
    field.go            # Field operations
    view.go             # View operations
    discussion.go       # Discussion operations
    workflow.go         # Workflow automation
    analytics.go        # Analytics and reporting

pkg/models/             # Shared data models
```

## Key Patterns

**Service Layer**: Commands in `internal/cmd/` delegate to services in `internal/service/` which use the GraphQL client in `internal/api/`.

**GraphQL Client**: `internal/api/client.go` wraps `shurcooL/graphql` with rate limiting (10 req/sec), automatic retry with exponential backoff, and OAuth2 authentication.

**Authentication**: Uses `gh auth token` from GitHub CLI first, falls back to `GITHUB_TOKEN` or `GH_TOKEN` environment variables.

**Configuration**: Uses Viper with config file at `~/.ghx.yaml`, environment variables prefixed with `GHX_`.

## Pre-commit Hooks

Lefthook is configured with strict pre-commit hooks that run sequentially:
1. `gofmt -s` and `goimports` on staged files
2. `go mod tidy` if go.mod changed
3. `golangci-lint run` (must pass)
4. `go test -short -race` (must pass)
5. Build verification

Commit messages must follow: `type(scope): description` where type is one of: feat, fix, docs, style, refactor, test, chore.

## Testing Notes

- Tests use `testify/assert` and `testify/require`
- Mock HTTP transport for GraphQL client tests: see `internal/api/client_test.go`
- Use `-short` flag to skip integration tests
