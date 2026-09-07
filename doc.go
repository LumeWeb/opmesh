// Package opmesh provides a typed, frontend-neutral operation registry with
// policy-enforced execution.
//
// # Registry
//
// Catalog (NewCatalog) is a concurrency-safe, instance-scoped registry of
// Operation descriptors. Operations are declared as plain data via
// NewOperation(OperationSpec{...}) with typed OperationArg inputs and a
// Handler that runs against the host application's services. Frontend adapters
// (CLI, MCP, UI) project their own command/tool surfaces from the same
// descriptors without duplicating operation metadata or policy logic.
//
// # Policy
//
// Invoke enforces the declared Interaction, Visibility, and Safety rules
// against the acting Actor — a model, the hosting app, or a human operator —
// and returns sentinel errors (ErrHumanRequired, ErrConfirmRequired, plus the
// declaration errors ErrDuplicateOperation/ErrUnknownOperation) that callers
// gate on with errors.Is. Destructive operations invoked by a model require
// explicit human confirmation unless the operation opted into agent
// self-confirmation via AgentConfirm on its `confirm` arg.
//
// # Normalization
//
// NormalizeOperationInput (also applied inside Invoke) coerces raw input into
// the declared arg shapes and fills declared defaults with a single
// resolveArg state machine, so every surface — flat flags, structured tool
// schemas, or direct Go calls — delivers identical, correctly-typed input to
// the Handler. It also enforces camelCase aliasing for kebab-case arg names,
// strict unknown-argument rejection, reserved plumbing keys
// (ReservedAuthTokenKey, ReservedRequestStateKey), and exactly-one-of
// SelectionGroup semantics (ErrSelector).
//
// # Helpers
//
// Typed input accessors (StrArg, IntArg, BoolArg, BoolArgPtr, IntArgPtr,
// StrSliceArg, StrFlexibleArg, SearchArg), canonical paging (ListArgs,
// ParseList, ParseListPage, ListOptions, PageLister, ScanPages), and
// positional mapping (MapPositionalArgs) round out the toolkit handlers use.
//
// # Boundaries
//
// opmesh is stdlib-only at runtime (testify is a test-only dependency) and
// imports no frontend frameworks (urfave/cli, MCP SDKs) or terminal
// renderers. It owns discovery semantics (Search/Describe, ToolDescriptor,
// the JSON-Schema builder) but not frontend presentation: CLI command
// compilation and MCP tool compilation are separate adapter concerns built on
// top of this registry.
package opmesh
