package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

func assertMatchingConfirmationHash(current values.ConfirmationHash, supplied values.ConfirmationHash) error {
	if !current.Equals(supplied) {
		return errors.Mark(errors.New("wrong confirmation hash supplied"), lib.ErrDomainConstraintsViolation)
	}

	return nil
}
