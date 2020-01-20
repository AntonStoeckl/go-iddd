package values

import (
	"go-iddd/service/lib"

	"github.com/cockroachdb/errors"
)

type PersonName struct {
	givenName  string
	familyName string
}

func BuildPersonName(givenName string, familyName string) (PersonName, error) {
	if familyName == "" {
		err := lib.MarkAndWrapError(
			errors.New("empty input for familyName"),
			lib.ErrInputIsInvalid,
			"BuildPersonName",
		)

		return PersonName{}, err
	}

	if givenName == "" {
		err := lib.MarkAndWrapError(
			errors.New("empty input for givenName"),
			lib.ErrInputIsInvalid,
			"BuildPersonName",
		)

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
