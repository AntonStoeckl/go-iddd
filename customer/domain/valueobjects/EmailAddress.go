package valueobjects

import (
	"encoding/json"
	"go-iddd/shared"
	"regexp"

	"golang.org/x/xerrors"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^[^\s]+@[^\s]+\.[\w]{2,}$`)
)

type EmailAddress struct {
	value string
}

/*** Factory methods ***/

func NewEmailAddress(from string) (*EmailAddress, error) {
	newEmailAddress := buildEmailAddress(from)

	if err := newEmailAddress.mustBeValid(); err != nil {
		return nil, err
	}

	return newEmailAddress, nil
}

func ReconstituteEmailAddress(from string) *EmailAddress {
	return buildEmailAddress(from)
}

func buildEmailAddress(from string) *EmailAddress {
	return &EmailAddress{value: from}
}

/*** Validation ***/

func (emailAddress *EmailAddress) mustBeValid() error {
	if matched := emailAddressRegExp.MatchString(emailAddress.value); matched != true {
		return xerrors.Errorf("emailAddress.mustBeValid: input does not match regex: %w", shared.ErrInvalidInput)
	}

	return nil
}

/*** Public methods implementing EmailAddress ***/

func (emailAddress *EmailAddress) EmailAddress() string {
	return emailAddress.value
}

func (emailAddress *EmailAddress) Equals(other *EmailAddress) bool {
	return emailAddress.EmailAddress() == other.EmailAddress()
}

func (emailAddress *EmailAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(emailAddress.value)
}

func UnmarshalEmailAddress(input interface{}) (*EmailAddress, error) {
	value, ok := input.(string)
	if !ok {
		return nil, xerrors.Errorf("UnmarshalEmailAddress: input is not a [string]: %w", shared.ErrUnmarshaling)
	}

	return &EmailAddress{value: value}, nil
}
