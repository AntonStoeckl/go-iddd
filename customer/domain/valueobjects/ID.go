package valueobjects

import (
	"encoding/json"
	"go-iddd/shared"

	"golang.org/x/xerrors"

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
	bytes, err := json.Marshal(id.value)
	if err != nil {
		return bytes, xerrors.Errorf("id.MarshalJSON: %s: %w", err, shared.ErrMarshaling)
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (id *ID) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return xerrors.Errorf("id.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshaling)
	}

	id.value = value

	return nil
}
