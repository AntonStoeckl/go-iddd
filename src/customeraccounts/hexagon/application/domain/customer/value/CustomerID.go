package value

import (
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type CustomerID string

func GenerateCustomerID() CustomerID {
	return CustomerID(uuid.New().String())
}

func BuildCustomerID(value string) (CustomerID, error) {
	if value == "" {
		err := errors.New("empty input for CustomerID")
		err = shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "BuildCustomerID")

		return "", err
	}

	id := CustomerID(value)

	return id, nil
}

func RebuildCustomerID(value string) CustomerID {
	return CustomerID(value)
}

func (id CustomerID) String() string {
	return string(id)
}

func (id CustomerID) Equals(other CustomerID) bool {
	return id.String() == other.String()
}
