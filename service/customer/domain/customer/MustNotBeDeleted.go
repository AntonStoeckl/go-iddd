package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

func MustNotBeDeleted(currentState state) error {
	if currentState.isDeleted {
		return errors.Mark(errors.New("customer is deleted"), lib.ErrNotFound)
	}

	return nil
}
