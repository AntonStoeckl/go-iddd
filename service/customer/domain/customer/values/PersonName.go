package values

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

type PersonName struct {
	givenName  string
	familyName string
}

func BuildPersonName(givenName string, familyName string) (PersonName, error) {
	wrapWithMsg := "BuildPersonName"

	if familyName == "" {
		err := errors.New("empty input for familyName")
		err = lib.MarkAndWrapError(err, lib.ErrInputIsInvalid, wrapWithMsg)

		return PersonName{}, err
	}

	if givenName == "" {
		err := errors.New("empty input for givenName")
		err = lib.MarkAndWrapError(err, lib.ErrInputIsInvalid, wrapWithMsg)

		return PersonName{}, err
	}

	personName := PersonName{
		givenName:  givenName,
		familyName: familyName,
	}

	return personName, nil
}

func RebuildPersonName(givenName string, familyName string) PersonName {
	personName := PersonName{
		givenName:  givenName,
		familyName: familyName,
	}

	return personName
}

func (personName PersonName) GivenName() string {
	return personName.givenName
}

func (personName PersonName) FamilyName() string {
	return personName.familyName
}

func (personName PersonName) Equals(other PersonName) bool {
	if personName.GivenName() != other.GivenName() {
		return false
	}

	if personName.FamilyName() != other.FamilyName() {
		return false
	}

	return true
}
