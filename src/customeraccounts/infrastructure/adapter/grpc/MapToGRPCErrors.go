package customergrpc

import (
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapToGRPCErrors(appErr error) error {
	var code codes.Code

	switch true {
	case errors.Is(appErr, shared.ErrInputIsInvalid):
		code = codes.InvalidArgument
	case errors.Is(appErr, shared.ErrNotFound):
		code = codes.NotFound
	case errors.Is(appErr, shared.ErrDuplicate):
		code = codes.AlreadyExists

	case errors.Is(appErr, shared.ErrDomainConstraintsViolation):
		code = codes.FailedPrecondition

	case errors.Is(appErr, shared.ErrMaxRetriesExceeded):
		code = codes.Aborted
	case errors.Is(appErr, shared.ErrConcurrencyConflict):
		code = codes.Aborted

	default:
		code = codes.Internal
	}

	return status.Errorf(code, "%s", errors.Cause(appErr))
}
