package valueobjects

import (
	"encoding/json"
	"errors"
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

func (personName *PersonName) MarshalJSON() ([]byte, error) {
	data := &struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	}{
		GivenName:  personName.givenName,
		FamilyName: personName.familyName,
	}

	return json.Marshal(data)
}

func UnmarshalPersonName(data interface{}) (*PersonName, error) {
	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("zefix")
	}

	personName := &PersonName{}

	for key, value := range values {
		value, ok := value.(string)
		if !ok {
			return nil, errors.New("zefix")
		}

		switch key {
		case "givenName":
			personName.givenName = value
		case "familyName":
			personName.familyName = value
		}
	}

	return personName, nil
}
