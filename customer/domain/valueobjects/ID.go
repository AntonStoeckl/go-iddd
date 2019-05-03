package valueobjects

import (
	"encoding/json"
	"errors"
	"go-iddd/shared"

	"github.com/google/uuid"
)

type ID interface {
	ID() string

	shared.AggregateIdentifier
}

type id struct {
	value string
}

/*** Factory methods ***/

func GenerateID() *id {
	return buildID(uuid.New().String())
}

func ReconstituteID(from string) *id {
	return buildID(from)
}

func buildID(from string) *id {
	return &id{value: from}
}

/*** Public methods implementing ID ***/

func (idenfifier *id) ID() string {
	return idenfifier.value
}

/*** Public methods implementing ID (methods for shared.AggregateIdentifier) ***/

func (idenfifier *id) String() string {
	return idenfifier.value
}

func (idenfifier *id) Equals(other shared.AggregateIdentifier) bool {
	if _, ok := other.(*id); !ok {
		return false
	}

	return idenfifier.String() == other.String()
}

func (idenfifier *id) MarshalJSON() ([]byte, error) {
	return json.Marshal(idenfifier.value)
}

func UnmarshalID(data interface{}) (*id, error) {
	value, ok := data.(string)
	if !ok {
		return nil, errors.New("zefix")
	}

	return &id{value: value}, nil
}
