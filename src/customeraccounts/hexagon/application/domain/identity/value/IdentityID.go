package value

import (
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type IdentityID string

func GenerateIdentityID() IdentityID {
	return IdentityID(uuid.New().String())
}

func BuildIdentityID(value string) (IdentityID, error) {
	if value == "" {
		err := errors.New("empty input for IdentityID")
		err = shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "BuildIdentityID")

		return "", err
	}

	id := IdentityID(value)

	return id, nil
}

func RebuildIdentityID(value string) IdentityID {
	return IdentityID(value)
}

func (id IdentityID) String() string {
	return string(id)
}

func (id IdentityID) Equals(other IdentityID) bool {
	return id.String() == other.String()
}
