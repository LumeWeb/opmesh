# AGENTS.md

This file provides development guidelines and architectural documentation for
the opmesh project.

## Common Commands

### Building
```bash
# Build all packages — opmesh is a pure Go library; no build tags (e.g.
# sqlite_fts5) are needed
go build ./...
```

### Testing
```bash
# Run all tests with race detection and coverage
go test -v -race -coverprofile=coverage.txt ./...

# Run tests for the root package only
go test -v -race .

# View coverage report
go tool cover -func=coverage.txt
```

### Mock Generation

```bash
# Generate mocks for interfaces (uses .mockery.yaml; mockery is
# pre-installed at $HOME/go/bin/mockery). Run with no args from the repo
# root. NEVER reinstall mockery.
mockery
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy dependencies
go mod tidy
```

## Project Overview

This project provides a typed, frontend-neutral operation registry with
policy-enforced execution, extracted from the Pinner CLI's operation
catalog (`internal/catalog`). Frontend compilers (CLI, MCP) build on
registries to project typed operations into their own command and tool
surfaces.

## Architecture

### Package Structure
- **Module path**: `go.lumeweb.com/opmesh`
- **Root package** (`opmesh`): a flat package containing all functionality;
  `doc.go` holds the package documentation
- **`mocks/`**: generated testify mocks of `opmesh` interfaces (mockery,
  configured via `.mockery.yaml`; committed, not ignored)
- **No `vendor/`**: dependencies are managed by the Go module system

### Design
- Registries hold typed operations; policy enforcement (human-in-the-loop,
  confirmation, actor visibility) happens at execution time and returns
  sentinel errors from `errors.go`, wrappable with `errors.Is`
- Frontend compilers consume registries; operation metadata and policy logic
  live here once, not per-frontend
- opmesh imports no frontend frameworks (urfave/cli, MCP SDKs)

### Testing Conventions
- Do not add a README/board copyright header to source files; attribution
  lives only in the LICENSE file
- Generated mocks live in `mocks/` (mockery, testify templates)
