package values

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type CustomerID struct {
	value string
}

func GenerateCustomerID() CustomerID {
	return CustomerID{value: uuid.New().String()}
}

func BuildCustomerID(value string) (CustomerID, error) {
	if value == "" {
		err := errors.New("empty input for CustomerID")
		err = lib.MarkAndWrapError(err, lib.ErrInputIsInvalid, "BuildCustomerID")

		return CustomerID{}, err
	}

	id := CustomerID{value: value}

	return id, nil
}

func RebuildCustomerID(value string) CustomerID {
	return CustomerID{value: value}
}

func (id CustomerID) String() string {
	return id.value
}

func (id CustomerID) Equals(other CustomerID) bool {
	return id.String() == other.String()
}
