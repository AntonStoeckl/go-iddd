package valueobjects

import (
	"encoding/json"
	"errors"
)

type PersonName interface {
	GivenName() string
	FamilyName() string
	Equals(other PersonName) bool
}

type personName struct {
	givenName  string
	familyName string
}

/*** Factory methods ***/

func NewPersonName(givenName string, familyName string) (*personName, error) {
	newPersonName := buildPersonName(givenName, familyName)
	if err := newPersonName.mustBeValid(); err != nil {
		return nil, err
	}

	return newPersonName, nil
}

func ReconstitutePersonName(givenName string, familyName string) *personName {
	return buildPersonName(givenName, familyName)
}

func buildPersonName(givenName string, familyName string) *personName {
	return &personName{
		givenName:  givenName,
		familyName: familyName,
	}
}

/*** Validation ***/

func (personName *personName) mustBeValid() error {
	if personName.familyName == "" {
		return errors.New("personName - empty input given for familyName")
	}

	if personName.givenName == "" {
		return errors.New("personName - empty input given for givenName")
	}

	return nil
}

/*** Public methods implementing PersonName ***/

func (personName *personName) GivenName() string {
	return personName.givenName
}

func (personName *personName) FamilyName() string {
	return personName.familyName
}

func (personName *personName) Equals(other PersonName) bool {
	if personName.GivenName() != other.GivenName() {
		return false
	}

	if personName.FamilyName() != other.FamilyName() {
		return false
	}

	return true
}

func (personName *personName) MarshalJSON() ([]byte, error) {
	data := &struct {
		GivenName  string `json:"givenName"`
		FamilyName string `json:"familyName"`
	}{
		GivenName:  personName.givenName,
		FamilyName: personName.familyName,
	}

	return json.Marshal(data)
}

func UnmarshalPersonName(data interface{}) (*personName, error) {
	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("zefix")
	}

	personName := &personName{}

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
