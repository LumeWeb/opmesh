package opmesh

import "errors"

// ErrDuplicateOperation is returned when an operation is registered under a
// name that already exists in the registry. Add wraps it with the conflicting
// operation name.
var ErrDuplicateOperation = errors.New("opmesh: operation already registered")

// ErrUnknownOperation is returned when an execution or lookup names an
// operation that is not registered. Invoke wraps it with the requested name.
var ErrUnknownOperation = errors.New("opmesh: unknown operation")

// ErrPolicyDenied is returned when policy enforcement refuses an execution —
// e.g. the operation requires a human or confirmation, or is not visible to
// the actor. It is the coarse umbrella: the specific policy sentinels below
// are the precise signals callers should gate on via errors.Is.
var ErrPolicyDenied = errors.New("opmesh: policy denied operation")

// ErrConfirmRequired is a sentinel signal that a destructive operation was
// invoked by a model actor and needs explicit human confirmation. Callers
// detect it with errors.Is(err, ErrConfirmRequired) to surface a confirm flow.
var ErrConfirmRequired = errors.New("confirmation required")

// ErrHumanRequired is a sentinel signal that an InteractionHumanOnly or
// InteractionNeedsHandoff operation was invoked by a non-human actor and so
// was refused at the gate. Such flows demand out-of-band human action a model
// or app would never complete. Callers detect it with
// errors.Is(err, ErrHumanRequired) to route the refusal to a human hand-off
// rather than treating it as a failure.
var ErrHumanRequired = errors.New("operation requires a human")

// ErrSelector is a sentinel signal that a SelectionGroup resolved to the wrong
// number of selected members. A SelectionGroup requires exactly one member to
// be selected; more-than-one is an ambiguous selector (e.g. a remove op given
// both its cids and all=true). Callers detect it with
// errors.Is(err, ErrSelector) to surface a selector contract violation.
var ErrSelector = errors.New("selector group must select exactly one member")
