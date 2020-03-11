package customer

import (
	"go-iddd/service/lib"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type ID struct {
	value string
}

func GenerateID() ID {
	return ID{value: uuid.New().String()}
}

func BuildID(value string) (ID, error) {
	if value == "" {
		err := lib.MarkAndWrapError(
			errors.New("empty input for ID"),
			lib.ErrInputIsInvalid,
			"BuildID",
		)

		return ID{}, err
	}

	id := ID{value: value}

	return id, nil
}

func RebuildID(value string) ID {
	return ID{value: value}
}

func (id ID) ID() string {
	return id.value
}

func (id ID) Equals(other ID) bool {
	return id.ID() == other.ID()
}
