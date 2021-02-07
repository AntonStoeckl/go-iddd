package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

func assertMatchingConfirmationHash(current, supplied value.ConfirmationHash) error {
	if !current.Equals(supplied) {
		return errors.Mark(errors.New("wrong confirmation hash supplied"), shared.ErrDomainConstraintsViolation)
	}

	return nil
}
