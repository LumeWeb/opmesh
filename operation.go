// Operation descriptors: the typed, frontend-neutral definitions of what can
// be registered and executed. All imports are stdlib-only by design.
package opmesh

import (
	"context"
	"encoding/json"
)

// Safety classifies what an operation does to state. It is declared on the
// operation, not inferred from a command name.
type Safety int

const (
	// SafetyRead is a pure read with no mutation.
	SafetyRead Safety = iota
	// SafetyMutate mutates state but is non-destructive.
	SafetyMutate
	// SafetyDestructive is irreversible (delete, forget, force ops).
	SafetyDestructive
)

// Interaction declares how a frontend may invoke the operation.
type Interaction int

const (
	// InteractionAgentSafe is fully automatable, no human input/OOB.
	InteractionAgentSafe Interaction = iota
	// InteractionHumanOnly prompts interactively; no agent-safe form.
	InteractionHumanOnly
	// InteractionNeedsHandoff requires external browser/device, split into 2 calls.
	InteractionNeedsHandoff
)

// Visibility declares which consumers may discover/invoke the operation.
type Visibility int

const (
	// VisibilityModel is agent-visible (searchable), the default.
	VisibilityModel Visibility = iota
	// VisibilityAppOnly is an app-only helper; invisible to agent discovery.
	VisibilityAppOnly
	// VisibilityBoth is both agent- and app-visible.
	VisibilityBoth
)

// ArgType classifies an operation argument's shape and drives both input
// normalization and the JSON-Schema property emitted by the schema builder.
type ArgType int

const (
	// ArgTypeString is a string argument.
	ArgTypeString ArgType = iota
	// ArgTypeBool is a boolean argument.
	ArgTypeBool
	// ArgTypeNullableBool is a tri-state boolean argument: its value may be
	// absent (nil), true, or false, and the Handler receives a *bool so it can
	// distinguish "omitted" from "explicitly false". Unlike ArgTypeBool — whose
	// absence is collapsed to false by normalizeInputDefaults — a nullable bool
	// preserves the absent state, which operations need when omission means
	// "leave unchanged" or "use the backend default" rather than "false".
	ArgTypeNullableBool
	// ArgTypeNullableInt is a tri-state integer argument: its value may be
	// absent (nil) or an explicit *int, and the Handler receives a *int so it
	// can distinguish "omitted" from "explicitly 0". Unlike ArgTypeInt — whose
	// absence is collapsed to 0 by normalizeInputDefaults — a nullable int
	// preserves the absent state, which operations need when omission means
	// "use the backend/op default" rather than "set 0" (e.g. a priority value
	// that defaults to 10 when omitted).
	ArgTypeNullableInt
	// ArgTypeInt is an integer argument.
	ArgTypeInt
	// ArgTypeFloat is a floating-point argument.
	ArgTypeFloat
	// ArgTypeDuration is a duration argument.
	ArgTypeDuration
	// ArgTypeStringSlice is a slice of strings argument.
	ArgTypeStringSlice
	// ArgTypeFlexibleID is a string argument that also accepts a JSON integer
	// (e.g. a list operation whose backend ids are ints). It coerce-renders any
	// accepted numeric to its decimal string form so a Handler reading via
	// StrArg/StrFlexibleArg always receives a string, and it advertises
	// ["string","integer"] in the JSON Schema so a caller can pass either the
	// id's integer form or a string form without being rejected by the
	// normalizer.
	ArgTypeFlexibleID
	// ArgTypeRawJSON is a pass-through argument: the value arrives as-is (a
	// JSON string or a decoded array/object) and the Handler is responsible for
	// parsing it. Use it with RawSchema to advertise a structured JSON Schema
	// on structured surfaces while letting flat surfaces treat it as a plain
	// string argument.
	ArgTypeRawJSON
)

// OperationArg describes one input. It drives input normalization/validation
// and the JSON Schema emitted by the schema builder; frontend adapters layer
// their own presentation (flags, env sources, positional bindings) on top
// without extra fields here.
type OperationArg struct {
	Name     string
	Type     ArgType
	Required bool
	// Default is the declared default in string form, parsed into the arg's
	// Go shape by the default filler. An arg that is Required but declares a
	// Default is satisfied by the default and so is optional from the
	// caller's perspective.
	Default string
	// Enum constrains a string arg to a closed set of values. Values are
	// matched case-insensitively at normalization time; the declared casing
	// is preserved in the value the Handler receives.
	Enum []string
	// Sensitive marks an arg whose value must be redacted (logs, echoed tool
	// calls) by any surface that renders it.
	Sensitive bool
	// Help is the single generic description of the arg, carried through the
	// JSON Schema as the property description. There is deliberately no
	// audience-specific variant: adapter packages that need tuned prose
	// compose it there.
	Help string
	// AgentRequired marks an arg required on the agent surface only. It is
	// never part of the generic required set used by NormalizeOperationInput,
	// so generic (non-agent) invocation paths do not enforce it. An adapter
	// that dispatches agent invocations applies its own requiredness check
	// over AgentRequired before handing input to the registry.
	AgentRequired bool
	// AgentConfirm marks a bool `confirm` arg as the explicit agent-confirmation
	// contract for a destructive operation: when the op is invoked by a model
	// actor and this arg's value is true, it satisfies the destructive safety
	// gate so the op runs headlessly without a human hand-off. It must
	// only be set on a bool arg literally named "confirm". Ops whose confirm
	// has a different semantic or that must always route through a human
	// hand-off must leave it unset so the model gate still refuses them with
	// ErrConfirmRequired.
	AgentConfirm bool
	// SelectionGroup groups mutually-exclusive selector members: exactly one
	// arg in a group may be selected per invocation (e.g. a "cids" slice and
	// an "all" bool). Empty selects membership in no group. Enforcement is
	// centralized in the normalize path so surfaces and direct Invoke all
	// agree.
	SelectionGroup string
	// RawSchema, when set, is used verbatim as the arg's JSON Schema property
	// instead of auto-generating one from ArgType. It lets an ArgTypeRawJSON
	// (or any) arg advertise a rich, closed structured schema that the flat
	// ArgType-to-JSON-Type mapping cannot express. RawSchema must be a JSON
	// object; it is merged as the property object (its "description" wins
	// over Help).
	RawSchema json.RawMessage
}

// Handler.Execute runs the business operation against the application's core
// services. It never touches frontend frameworks. Implemented by service code.
type Handler interface {
	Execute(ctx context.Context, input map[string]any) (any, error)
}

// Operation is the canonical descriptor consumed by every frontend. It is an
// interface so concrete operations can be different types; the registry holds
// them polymorphically.
type Operation interface {
	Name() string
	Title() string
	Summary() string
	Description() string
	Args() []OperationArg
	Positional() string
	Safety() Safety
	Interaction() Interaction
	Visibility() Visibility
	Category() string
	Handler() Handler
}

// OperationSpec is the plain struct NewOperation turns into an Operation.
type OperationSpec struct {
	Name        string
	Title       string
	Summary     string
	Description string
	Args        []OperationArg
	Positional  string
	Safety      Safety
	Interaction Interaction
	Visibility  Visibility
	Category    string
	Handler     Handler
}

// NewOperation builds a concrete Operation from a spec. Most operations use
// this; a bespoke operation may implement Operation directly (e.g. one whose
// Handler needs closures over core deps).
func NewOperation(spec OperationSpec) Operation { return simpleOperation{spec} }

// simpleOperation is the default Operation implementation, delegating every
// method to the fields of the spec it was constructed from.
type simpleOperation struct{ spec OperationSpec }

func (o simpleOperation) Name() string             { return o.spec.Name }
func (o simpleOperation) Title() string            { return o.spec.Title }
func (o simpleOperation) Summary() string          { return o.spec.Summary }
func (o simpleOperation) Description() string      { return o.spec.Description }
func (o simpleOperation) Args() []OperationArg     { return o.spec.Args }
func (o simpleOperation) Positional() string       { return o.spec.Positional }
func (o simpleOperation) Safety() Safety           { return o.spec.Safety }
func (o simpleOperation) Interaction() Interaction { return o.spec.Interaction }
func (o simpleOperation) Visibility() Visibility   { return o.spec.Visibility }
func (o simpleOperation) Category() string         { return o.spec.Category }
func (o simpleOperation) Handler() Handler         { return o.spec.Handler }
