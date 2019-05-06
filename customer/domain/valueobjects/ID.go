package valueobjects

import (
	"encoding/json"
	"errors"
	"go-iddd/shared"

	"github.com/google/uuid"
)

type ID struct {
	value string
}

/*** Factory methods ***/

func GenerateID() *ID {
	return buildID(uuid.New().String())
}

func ReconstituteID(from string) *ID {
	return buildID(from)
}

func buildID(from string) *ID {
	return &ID{value: from}
}

/*** Getter methods ***/

func (id *ID) ID() string {
	return id.value
}

/*** Implement shared.AggregateIdentifier ***/

func (id *ID) String() string {
	return id.value
}

func (id *ID) Equals(other shared.AggregateIdentifier) bool {
	if _, ok := other.(*ID); !ok {
		return false
	}

	return id.String() == other.String()
}

/*** Implement json.Marshaler ***/

func (id *ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.value)
}

/*** Another factory method ***/

func UnmarshalID(data interface{}) (*ID, error) {
	value, ok := data.(string)
	if !ok {
		return nil, errors.New("zefix")
	}

	return &ID{value: value}, nil
}
