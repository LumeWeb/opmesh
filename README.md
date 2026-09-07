opmesh
======

Go package providing a typed, frontend-neutral operation registry with
policy-enforced execution. Runtime stdlib-only (testify is a test-only
dependency); imports no frontend frameworks (urfave/cli, MCP SDKs) or
renderers. Frontend compilers (CLI, MCP) are separate adapters that build on
this registry.

Usage
-----

```go
import "go.lumeweb.com/opmesh"
```

Declare a typed operation and register it:

```go
op := opmesh.NewOperation(opmesh.OperationSpec{
    Name:        "pin.add",
    Title:       "Add Pin",
    Summary:     "Pin a CID",
    Category:    "pinning",
    Safety:      opmesh.SafetyRead,
    Interaction: opmesh.InteractionAgentSafe,
    Visibility:  opmesh.VisibilityModel,
    Args: []opmesh.OperationArg{
        {Name: "cid", Type: opmesh.ArgTypeString, Required: true, Help: "content id"},
 
    },
    Handler: myHandler{}, // opmesh.Handler: Execute(ctx, map[string]any) (any, error)
})

reg := opmesh.NewCatalog()
if err := reg.Add(op); err != nil {
    return err // wraps opmesh.ErrDuplicateOperation; check with errors.Is
}
```

Discover and execute with policy enforcement:

```go
desc, ok := reg.Describe("pin.add", opmesh.ActorModel) // discovery view + JSON Schema

res, err := reg.Invoke(ctx, "pin.add", map[string]any{"cid": "QmX"}, opmesh.ActorModel)
if err != nil {
    switch {
    case errors.Is(err, opmesh.ErrConfirmRequired): // destructive op needs confirmation
    case errors.Is(err, opmesh.ErrHumanRequired):   // human-only op denied for a model
    case errors.Is(err, opmesh.ErrUnknownOperation): // no such operation
    }
}
```

Normalize input for a Handler invoked directly (same coercion, defaults, and
selection-group gates as `Invoke`):

```go
normalized, err := opmesh.NormalizeOperationInput(op, input)
```

Development
-----------

```sh
go build ./...
go test -v -race -coverprofile=coverage.txt ./...
mockery   # when .mockery.yaml declares interfaces
```

## License

MIT — see [LICENSE](LICENSE).
