package shared

import "github.com/cockroachdb/errors"

var (
	ErrInputIsInvalid             = errors.New("input is invalid")
	ErrMarshalingFailed           = errors.New("marshaling failed")
	ErrUnmarshalingFailed         = errors.New("unmarshaling failed")
	ErrCommandIsInvalid           = errors.New("command is invalid")
	ErrCommandIsUnknown           = errors.New("command is unknown")
	ErrDomainConstraintsViolation = errors.New("domain constraints violation")
	ErrConcurrencyConflict        = errors.New("concurrency conflict")
	ErrMaxRetriesExceeded         = errors.New("max retries exceeded")
	ErrNotFound                   = errors.New("not found")
	ErrDuplicate                  = errors.New("duplicate")
	ErrTechnical                  = errors.New("technical")
)

func MarkAndWrapError(original error, markAs error, wrapWith string) error {
	return errors.Mark(errors.Wrap(original, wrapWith), markAs)
}
