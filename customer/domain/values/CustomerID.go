package values

import (
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type CustomerID struct {
	value string
}

/*** Factory methods ***/

func GenerateCustomerID() *CustomerID {
	return buildCustomerID(uuid.New().String())
}

func RebuildCustomerID(from string) (*CustomerID, error) {
	rebuiltID := buildCustomerID(from)

	if err := rebuiltID.shouldBeValid(); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInputIsInvalid), "id.New")
	}

	return rebuiltID, nil
}

func buildCustomerID(from string) *CustomerID {
	return &CustomerID{value: from}
}

func (id *CustomerID) shouldBeValid() error {
	if id.value == "" {
		return errors.New("empty input for CustomerID")
	}

	return nil
}

/*** Getter Methods (implement shared.IdentifiesAggregates) ***/

func (id *CustomerID) String() string {
	return id.value
}

/*** Comparison Methods (implement shared.IdentifiesAggregates) ***/

func (id *CustomerID) Equals(other shared.IdentifiesAggregates) bool {
	if _, ok := other.(*CustomerID); !ok {
		return false
	}

	return id.String() == other.String()
}

/*** Implement json.Marshaler ***/

func (id *CustomerID) MarshalJSON() ([]byte, error) {
	bytes, err := jsoniter.Marshal(id.value)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrMarshalingFailed), "CustomerID.MarshalJSON")
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (id *CustomerID) UnmarshalJSON(data []byte) error {
	var value string

	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "CustomerID.UnmarshalJSON")
	}

	id.value = value

	return nil
}
