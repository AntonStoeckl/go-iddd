package valueobjects

import (
	"go-iddd/shared"

	"github.com/google/uuid"
)

type ID interface {
	shared.AggregateIdentifier
}

type id struct {
	value string
}

func GenerateID() *id {
	return ReconstituteID(uuid.New().String())
}

func ReconstituteID(from string) *id {
	return &id{value: from}
}

func (idenfifier *id) String() string {
	return idenfifier.value
}

func (idenfifier *id) Equals(other shared.AggregateIdentifier) bool {
	if _, ok := other.(*id); !ok {
		return false
	}

	return idenfifier.String() == other.String()
}
