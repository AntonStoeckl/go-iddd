package valueobjects

import (
	"encoding/json"
	"errors"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

type PersonName struct {
	givenName  string
	familyName string
}

/*** Factory methods ***/

func NewPersonName(givenName string, familyName string) (*PersonName, error) {
	newPersonName := buildPersonName(givenName, familyName)
	if err := newPersonName.mustBeValid(); err != nil {
		return nil, err
	}

	return newPersonName, nil
}

func ReconstitutePersonName(givenName string, familyName string) *PersonName {
	return buildPersonName(givenName, familyName)
}

func buildPersonName(givenName string, familyName string) *PersonName {
	return &PersonName{
		givenName:  givenName,
		familyName: familyName,
	}
}

/*** Validation ***/

func (personName *PersonName) mustBeValid() error {
	if personName.familyName == "" {
		return errors.New("personName - empty input given for familyName")
	}

	if personName.givenName == "" {
		return errors.New("personName - empty input given for givenName")
	}

	return nil
}

/*** Public methods implementing PersonName ***/

func (personName *PersonName) GivenName() string {
	return personName.givenName
}

func (personName *PersonName) FamilyName() string {
	return personName.familyName
}

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

	bytes, err := json.Marshal(data)
	if err != nil {
		return bytes, xerrors.Errorf("personName.MarshalJSON: %s: %w", err, shared.ErrMarshaling)
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (personName *PersonName) UnmarshalJSON(data []byte) error {
	values := &struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	}{}

	if err := json.Unmarshal(data, values); err != nil {
		return xerrors.Errorf("personName.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshaling)
	}

	personName.givenName = values.GivenName
	personName.familyName = values.FamilyName

	return nil
}
