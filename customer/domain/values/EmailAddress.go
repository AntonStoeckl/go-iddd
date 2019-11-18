package values

import (
	"go-iddd/shared"
	"regexp"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^[^\s]+@[^\s]+\.[\w]{2,}$`)
)

type EmailAddress struct {
	value string
}

/*** Factory methods ***/

func EmailAddressFrom(input string) (*EmailAddress, error) {
	newEmailAddress := &EmailAddress{value: input}

	if err := newEmailAddress.shouldBeValid(); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInputIsInvalid), "emailAddress.New")
	}

	return newEmailAddress, nil
}

func (emailAddress *EmailAddress) shouldBeValid() error {
	if matched := emailAddressRegExp.MatchString(emailAddress.value); matched != true {
		return errors.New("input has invalid format")
	}

	return nil
}

/*** Getter Methods ***/

func (emailAddress *EmailAddress) EmailAddress() string {
	return emailAddress.value
}

/*** Comparison Methods ***/

func (emailAddress *EmailAddress) Equals(other *EmailAddress) bool {
	return emailAddress.value == other.value
}

func (emailAddress *EmailAddress) ShouldEqual(other *EmailAddress) error {
	if !emailAddress.Equals(other) {
		return errors.Mark(errors.New("emailAddress.ShouldEqual"), shared.ErrNotEqual)
	}

	return nil
}

/*** Conversion Methods ***/

func (emailAddress *EmailAddress) ToConfirmable() *ConfirmableEmailAddress {
	return &ConfirmableEmailAddress{
		emailAddress:     emailAddress,
		confirmationHash: GenerateConfirmationHash(emailAddress.EmailAddress()),
		isConfirmed:      false,
	}
}

/*** Implement json.Marshaler ***/

func (emailAddress *EmailAddress) MarshalJSON() ([]byte, error) {
	bytes, err := jsoniter.Marshal(emailAddress.value)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrMarshalingFailed), "emailAddress.MarshalJSON")
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (emailAddress *EmailAddress) UnmarshalJSON(data []byte) error {
	var value string

	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "emailAddress.UnmarshalJSON")
	}

	emailAddress.value = value

	return nil
}
