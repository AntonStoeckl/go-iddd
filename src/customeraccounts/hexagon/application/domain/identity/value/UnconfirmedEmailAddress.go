package value

import (
	"regexp"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^\S+@\S+\.\w{2,}$`)
)

type UnconfirmedEmailAddress struct {
	value            string
	confirmationHash ConfirmationHash
}

func BuildUnconfirmedEmailAddress(input string) (UnconfirmedEmailAddress, error) {
	if matched := emailAddressRegExp.MatchString(input); !matched {
		err := errors.New("input has invalid format")
		err = shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "UnconfirmedEmailAddress")

		return UnconfirmedEmailAddress{}, err
	}

	emailAddress := UnconfirmedEmailAddress{
		value:            input,
		confirmationHash: GenerateConfirmationHash(input),
	}

	return emailAddress, nil
}

func RebuildUnconfirmedEmailAddress(input, hash string) UnconfirmedEmailAddress {
	return UnconfirmedEmailAddress{
		value:            input,
		confirmationHash: RebuildConfirmationHash(hash),
	}
}

func (emailAddress UnconfirmedEmailAddress) String() string {
	return emailAddress.value
}

func (emailAddress UnconfirmedEmailAddress) ConfirmationHash() ConfirmationHash {
	return emailAddress.confirmationHash
}

func (emailAddress UnconfirmedEmailAddress) Equals(other EmailAddress) bool {
	return emailAddress.String() == other.String()
}
