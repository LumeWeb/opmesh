package opmesh

import "errors"

// ErrDuplicateOperation is returned when an operation is registered under a
// name that already exists in the registry.
var ErrDuplicateOperation = errors.New("opmesh: operation already registered")

// ErrUnknownOperation is returned when an execution or lookup names an
// operation that is not registered.
var ErrUnknownOperation = errors.New("opmesh: unknown operation")

// ErrPolicyDenied is returned when policy enforcement refuses an execution —
// e.g. the operation requires a human or confirmation, or is not visible to
// the actor.
var ErrPolicyDenied = errors.New("opmesh: policy denied operation")
