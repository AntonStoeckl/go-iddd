package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

var wasDeletedErr = errors.Mark(errors.New("customer was deleted"), lib.ErrNotFound)

func wasNotDeleted(currentState currentState) bool {
	return !currentState.isDeleted
}
