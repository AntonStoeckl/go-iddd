package shared

import "errors"

var (
	ErrInputIsInvalid             = errors.New("input is invalid")
	ErrNotEqual                   = errors.New("not equal")
	ErrMarshalingFailed           = errors.New("marshaling failed")
	ErrUnmarshalingFailed         = errors.New("unmarshaling failed")
	ErrCommandIsInvalid           = errors.New("command is invalid")
	ErrCommandIsUnknown           = errors.New("command is unknown")
	ErrDomainConstraintsViolation = errors.New("domain constraints violation")
	ErrInvalidEventStream         = errors.New("invalid event stream")
	ErrConcurrencyConflict        = errors.New("concurrency conflict")
	ErrMaxRetriesExceeded         = errors.New("max retries exceeded")
	ErrNotFound                   = errors.New("not found")
	ErrDuplicate                  = errors.New("duplicate")
	ErrTechnical                  = errors.New("technical")
)
