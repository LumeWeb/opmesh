# opmesh

Go package providing a typed, frontend-neutral operation registry with
policy-enforced execution, extracted from the Pinner CLI's operation catalog.
Frontend compilers (CLI, MCP) build on registries.

## Usage

```go
import "go.lumeweb.com/opmesh"
```

Construct a registry and register a typed operation:

```go
reg, err := opmesh.New()
if err != nil {
    return err
}

if err := reg.Register(myOperation); err != nil {
    return err // wraps opmesh.ErrDuplicateOperation if the name is taken
}
```

Enforce policy and execute:

```go
res, err := reg.Execute(ctx, myOperation.Name(), opmesh.WithActor(actor), args)
if err != nil {
    switch {
    case errors.Is(err, opmesh.ErrPolicyDenied):
        // operation requires a human or confirmation, or is not visible
    case errors.Is(err, opmesh.ErrUnknownOperation):
        // no such operation
    }
}
```

## Development

```sh
go build ./...
go test -v -race -coverprofile=coverage.txt ./...
mockery
```

## License

MIT — see [LICENSE](LICENSE).
