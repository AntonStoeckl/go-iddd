package valueobjects

import (
    "errors"
    "regexp"
)

var (
    emailAddressRegExp = regexp.MustCompile(`^[^\s]+@[^\s]+\.[\w]{2,}$`)
)

type EmailAddress interface {
    String() string
    Equals(other EmailAddress) bool
}

type emailAddress struct {
    value string
}

func NewEmailAddress(from string) (*emailAddress, error) {
    newEmailAddress := ReconstituteEmailAddress(from)

    if err := newEmailAddress.mustBeValid(); err != nil {
        return nil, err
    }

    return newEmailAddress, nil
}

func (emailAddress *emailAddress) mustBeValid() error {
    if matched := emailAddressRegExp.MatchString(emailAddress.value); matched != true {
        return errors.New("emailAddress - invalid input given")
    }

    return nil
}

func ReconstituteEmailAddress(from string) *emailAddress {
    return &emailAddress{value: from}
}

func (emailAddress *emailAddress) String() string {
    return emailAddress.value
}

func (emailAddress *emailAddress) Equals(other EmailAddress) bool {
    return emailAddress.String() == other.String()
}
