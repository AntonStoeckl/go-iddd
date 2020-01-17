package values

import (
	"go-iddd/shared"

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
		err := shared.MarkAndWrapError(
			errors.New("empty input for CustomerID"),
			shared.ErrInputIsInvalid,
			"BuildCustomerID",
		)

		return CustomerID{}, err
	}

	id := CustomerID{value: value}

	return id, nil
}

func RebuildCustomerID(value string) CustomerID {
	return CustomerID{value: value}
}

func (id CustomerID) ID() string {
	return id.value
}

func (id CustomerID) String() string {
	return id.value
}

func (id CustomerID) Equals(other shared.IdentifiesAggregates) bool {
	return id.String() == other.String()
}
