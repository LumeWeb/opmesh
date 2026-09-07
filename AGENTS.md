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
policy-enforced execution. Runtime stdlib-only (testify is a test-only
dependency) — opmesh imports no frontend frameworks (urfave/cli, MCP SDKs)
or renderers. Frontend compilers (CLI, MCP) are separate adapter packages
that project operations from this registry into their own surfaces.

## Architecture

### Package Structure
- **Module path**: `go.lumeweb.com/opmesh`
- **Root package** (`opmesh`): a flat package containing all functionality;
  `doc.go` holds the package documentation
- **`mocks/`**: generated testify mocks of `opmesh` interfaces (mockery,
  configured via `.mockery.yaml`; committed, not ignored; no mocks generated
  while `interfaces:` in `.mockery.yaml` is empty)
- **No `vendor/`**: dependencies are managed by the Go module system

### Design
- `Catalog` (`NewCatalog`) is an instance-scoped, concurrency-safe registry
  of `Operation` descriptors (`NewOperation(OperationSpec{...})`); it owns
  discovery (`Search`/`Get`/`Describe` → `ToolDescriptor` + JSON Schema) and
  dispatch (`Invoke`)
- Policy enforcement (interaction human-in-the-loop, actor visibility,
  destructive confirmation) happens at dispatch time against the acting
  `Actor` and returns sentinel errors from `errors.go`
  (`ErrConfirmRequired`, `ErrHumanRequired`, `ErrSelector`,
  `ErrDuplicateOperation`, `ErrUnknownOperation`), checkable with `errors.Is`
- Input normalization (`NormalizeOperationInput`) coerces raw values to
  declared arg shapes, fills declared defaults, aliases camelCase → kebab
  names, rejects unknown args, keeps reserved keys
  (`ReservedAuthTokenKey`/`ReservedRequestStateKey`), and enforces
  exactly-one-of `SelectionGroup` semantics (`ErrSelector`)
- Helpers for handlers: typed accessors (`StrArg`, `IntArg`, `BoolArgPtr`,
  `IntArgPtr`, `StrSliceArg`, `StrFlexibleArg`, `SearchArg`), paging
  (`ListArgs`, `ParseList`, `ParseListPage`, `ListOptions`, `PageLister`,
  `ScanPages`), positional mapping (`MapPositionalArgs`)
- No mutable package-global state: registry state is instance-scoped;
  reserved keys are constants
- What opmesh does NOT contain: CLI flag/command compilation, MCP tool
  compilation, environment-source handling, or audience-tuned descriptions —
  all of that belongs to frontend adapter layers

### Testing Conventions
- Do not add a README/board copyright header to source files; attribution
  lives only in the LICENSE file
- Generated mocks live in `mocks/` (mockery, testify templates)
