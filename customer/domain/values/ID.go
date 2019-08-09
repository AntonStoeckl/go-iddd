package values

import (
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type ID struct {
	value string
}

/*** Factory methods ***/

func GenerateID() *ID {
	return buildID(uuid.New().String())
}

func RebuildID(from string) (*ID, error) {
	rebuiltID := buildID(from)

	if err := rebuiltID.shouldBeValid(); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInputIsInvalid), "id.New")
	}

	return rebuiltID, nil
}

func buildID(from string) *ID {
	return &ID{value: from}
}

func (id *ID) shouldBeValid() error {
	if id.value == "" {
		return errors.New("empty input for id")
	}

	return nil
}

/*** Getter Methods (implement shared.IdentifiesAggregates) ***/

func (id *ID) String() string {
	return id.value
}

/*** Comparison Methods (implement shared.IdentifiesAggregates) ***/

func (id *ID) Equals(other shared.IdentifiesAggregates) bool {
	if _, ok := other.(*ID); !ok {
		return false
	}

	return id.String() == other.String()
}

/*** Implement json.Marshaler ***/

func (id *ID) MarshalJSON() ([]byte, error) {
	bytes, err := jsoniter.Marshal(id.value)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrMarshalingFailed), "id.MarshalJSON")
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (id *ID) UnmarshalJSON(data []byte) error {
	var value string

	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "id.UnmarshalJSON")
	}

	id.value = value

	return nil
}
