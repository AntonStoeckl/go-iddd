package shared

import "errors"

var (
	ErrInputIsInvalid             = errors.New("input is invalid")
	ErrNotEqual                   = errors.New("not equal")
	ErrMarshalingFailed           = errors.New("marshaling failed")
	ErrUnmarshalingFailed         = errors.New("unmarshaling failed")
	ErrCommandCanNotBeHandled     = errors.New("command can not be handled")
	ErrCommandIsInvalid           = errors.New("command is invalid")
	ErrDomainConstraintsViolation = errors.New("domain constraints violation")
	ErrInvalidEventStream         = errors.New("invalid event stream")
	ErrConcurrencyConflict        = errors.New("concurrency conflict")
	ErrNotFound                   = errors.New("not found")
	ErrDuplicate                  = errors.New("duplicate")
)
