package values

import (
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

type PersonName struct {
	givenName  string
	familyName string
}

/*** Factory methods ***/

func NewPersonName(givenName string, familyName string) (*PersonName, error) {
	newPersonName := buildPersonName(givenName, familyName)

	if err := newPersonName.shouldBeValid(); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInputIsInvalid), "personName.New")
	}

	return newPersonName, nil
}

func buildPersonName(givenName string, familyName string) *PersonName {
	return &PersonName{
		givenName:  givenName,
		familyName: familyName,
	}
}

func (personName *PersonName) shouldBeValid() error {
	if personName.familyName == "" {
		return errors.New("empty input for familyName")
	}

	if personName.givenName == "" {
		return errors.New("empty input for givenName")
	}

	return nil
}

/*** Getter methods ***/

func (personName *PersonName) GivenName() string {
	return personName.givenName
}

func (personName *PersonName) FamilyName() string {
	return personName.familyName
}

/*** Comparison methods ***/

func (personName *PersonName) Equals(other *PersonName) bool {
	if personName.GivenName() != other.GivenName() {
		return false
	}

	if personName.FamilyName() != other.FamilyName() {
		return false
	}

	return true
}

/*** Implement json.Marshaler ***/

func (personName *PersonName) MarshalJSON() ([]byte, error) {
	data := &struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	}{
		GivenName:  personName.givenName,
		FamilyName: personName.familyName,
	}

	bytes, err := jsoniter.Marshal(data)
	if err != nil {
		return bytes, errors.Wrap(errors.Mark(err, shared.ErrMarshalingFailed), "personName.MarshalJSON")
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (personName *PersonName) UnmarshalJSON(data []byte) error {
	values := &struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	}{}

	if err := jsoniter.Unmarshal(data, values); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "personName.UnmarshalJSON")
	}

	personName.givenName = values.GivenName
	personName.familyName = values.FamilyName

	return nil
}
