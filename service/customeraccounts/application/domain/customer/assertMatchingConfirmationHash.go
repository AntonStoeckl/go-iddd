package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/cockroachdb/errors"
)

func assertMatchingConfirmationHash(current value.ConfirmationHash, supplied value.ConfirmationHash) error {
	if !current.Equals(supplied) {
		return errors.Mark(errors.New("wrong confirmation hash supplied"), shared.ErrDomainConstraintsViolation)
	}

	return nil
}
