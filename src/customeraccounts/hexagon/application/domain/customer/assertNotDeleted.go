package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

func assertNotDeleted(currentState currentState) error {
	if currentState.isDeleted {
		return errors.Mark(errors.New("customer was deleted"), shared.ErrNotFound)
	}

	return nil
}
