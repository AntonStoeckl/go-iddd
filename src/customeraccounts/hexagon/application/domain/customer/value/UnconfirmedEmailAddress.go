package value

import (
	"regexp"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^\S+@\S+\.\w{2,}$`)
)

type UnconfirmedEmailAddress string

func BuildUnconfirmedEmailAddress(input string) (UnconfirmedEmailAddress, error) {
	if matched := emailAddressRegExp.MatchString(input); !matched {
		err := errors.New("input has invalid format")
		err = shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "UnconfirmedEmailAddress")

		return "", err
	}

	emailAddress := UnconfirmedEmailAddress(input)

	return emailAddress, nil
}

func RebuildUnconfirmedEmailAddress(input string) UnconfirmedEmailAddress {
	return UnconfirmedEmailAddress(input)
}

func (emailAddress UnconfirmedEmailAddress) String() string {
	return string(emailAddress)
}

func (emailAddress UnconfirmedEmailAddress) Equals(other EmailAddress) bool {
	return emailAddress.String() == other.String()
}
