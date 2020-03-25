package customergrpc

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapToGRPCErrors(appErr error) error {
	var code codes.Code

	switch true {
	case errors.Is(appErr, lib.ErrInputIsInvalid):
		code = codes.InvalidArgument
	case errors.Is(appErr, lib.ErrNotFound):
		code = codes.NotFound
	case errors.Is(appErr, lib.ErrDuplicate):
		code = codes.AlreadyExists

	case errors.Is(appErr, lib.ErrDomainConstraintsViolation):
		code = codes.FailedPrecondition

	case errors.Is(appErr, lib.ErrMaxRetriesExceeded):
		code = codes.Aborted
	case errors.Is(appErr, lib.ErrConcurrencyConflict):
		code = codes.Aborted

	default:
		code = codes.Internal
	}

	return status.Errorf(code, "%s", errors.Cause(appErr))
}
