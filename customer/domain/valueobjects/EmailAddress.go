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

/*** Implement json.Marshaler ***/

func (emailAddress *EmailAddress) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(emailAddress.value)
	if err != nil {
		return nil, xerrors.Errorf("emailAddress.MarshalJSON: %s: %w", err, shared.ErrMarshaling)
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (emailAddress *EmailAddress) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return xerrors.Errorf("emailAddress.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshaling)
	}

	emailAddress.value = value

	return nil
}
