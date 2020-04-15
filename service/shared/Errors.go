package shared

import "github.com/cockroachdb/errors"

var (
	ErrInputIsInvalid = errors.New("input is invalid")
	ErrNotFound       = errors.New("not found")
	ErrDuplicate      = errors.New("duplicate")

	ErrDomainConstraintsViolation = errors.New("domain constraints violation")

	ErrMaxRetriesExceeded  = errors.New("max retries exceeded")
	ErrConcurrencyConflict = errors.New("concurrency conflict")

	ErrMarshalingFailed   = errors.New("marshaling failed")
	ErrUnmarshalingFailed = errors.New("unmarshaling failed")
	ErrTechnical          = errors.New("technical")
)

func MarkAndWrapError(original error, markAs error, wrapWith string) error {
	return errors.Mark(errors.Wrap(original, wrapWith), markAs)
}
