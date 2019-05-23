package values

import (
	"encoding/json"
	"errors"
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

	if err := newEmailAddress.shouldBeValid(); err != nil {
		return nil, xerrors.Errorf("emailAddress.New -> %s: %w", err, shared.ErrInputIsInvalid)
	}

	return newEmailAddress, nil
}

func (emailAddress *EmailAddress) shouldBeValid() error {
	if matched := emailAddressRegExp.MatchString(emailAddress.value); matched != true {
		return errors.New("input does not match regex")
	}

	return nil
}

func buildEmailAddress(from string) *EmailAddress {
	return &EmailAddress{value: from}
}

/*** Getter Methods ***/

func (emailAddress *EmailAddress) EmailAddress() string {
	return emailAddress.value
}

/*** Comparison Methods ***/

func (emailAddress *EmailAddress) Equals(other *EmailAddress) bool {
	return emailAddress.value == other.value
}

/*** Conversion Methods ***/

func (emailAddress *EmailAddress) ToConfirmable() *ConfirmableEmailAddress {
	return buildConfirmableEmailAddress(emailAddress, GenerateConfirmationHash(emailAddress.EmailAddress()))
}

/*** Implement json.Marshaler ***/

func (emailAddress *EmailAddress) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(emailAddress.value)
	if err != nil {
		return nil, xerrors.Errorf("emailAddress.MarshalJSON -> %s: %w", err, shared.ErrMarshalingFailed)
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (emailAddress *EmailAddress) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return xerrors.Errorf("emailAddress.UnmarshalJSON -> %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	emailAddress.value = value

	return nil
}
