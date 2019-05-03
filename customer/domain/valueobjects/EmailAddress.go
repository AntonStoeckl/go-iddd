package valueobjects

import (
	"encoding/json"
	"errors"
	"regexp"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^[^\s]+@[^\s]+\.[\w]{2,}$`)
)

type EmailAddress interface {
	EmailAddress() string
	Equals(other EmailAddress) bool
}

type emailAddress struct {
	value string
}

/*** Factory methods ***/

func NewEmailAddress(from string) (*emailAddress, error) {
	newEmailAddress := buildEmailAddress(from)

	if err := newEmailAddress.mustBeValid(); err != nil {
		return nil, err
	}

	return newEmailAddress, nil
}

func ReconstituteEmailAddress(from string) *emailAddress {
	return buildEmailAddress(from)
}

func buildEmailAddress(from string) *emailAddress {
	return &emailAddress{value: from}
}

/*** Validation ***/

func (emailAddress *emailAddress) mustBeValid() error {
	if matched := emailAddressRegExp.MatchString(emailAddress.value); matched != true {
		return errors.New("emailAddress - invalid input given")
	}

	return nil
}

/*** Public methods implementing EmailAddress ***/

func (emailAddress *emailAddress) EmailAddress() string {
	return emailAddress.value
}

func (emailAddress *emailAddress) Equals(other EmailAddress) bool {
	return emailAddress.EmailAddress() == other.EmailAddress()
}

func (emailAddress *emailAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(emailAddress.value)
}

func UnmarshalEmailAddress(data interface{}) (*emailAddress, error) {
	value, ok := data.(string)
	if !ok {
		return nil, errors.New("zefix")
	}

	return &emailAddress{value: value}, nil
}
