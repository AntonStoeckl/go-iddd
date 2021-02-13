package value

import (
	"regexp"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^\S+@\S+\.\w{2,}$`)
)

type EmailAddress string

func BuildEmailAddress(input string) (EmailAddress, error) {
	if matched := emailAddressRegExp.MatchString(input); !matched {
		err := errors.New("input has invalid format")
		err = shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "BuildEmailAddress")

		return "", err
	}

	emailAddress := EmailAddress(input)

	return emailAddress, nil
}

func RebuildEmailAddress(input string) EmailAddress {
	return EmailAddress(input)
}

func (emailAddress EmailAddress) String() string {
	return string(emailAddress)
}

func (emailAddress EmailAddress) Equals(other EmailAddress) bool {
	return emailAddress.String() == other.String()
}
