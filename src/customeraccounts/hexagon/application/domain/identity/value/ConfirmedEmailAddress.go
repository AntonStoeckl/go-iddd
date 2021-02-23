package value

import (
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

type ConfirmedEmailAddress string

func ConfirmEmailAddressWithHash(
	emailAddress UnconfirmedEmailAddress,
	confirmationHash ConfirmationHash,
) (ConfirmedEmailAddress, error) {

	if !emailAddress.confirmationHash.Equals(confirmationHash) {
		return "", errors.Mark(
			errors.New("confirmEmailAddressWithHash: wrong confirmation hash supplied"),
			shared.ErrDomainConstraintsViolation,
		)
	}

	return ConfirmedEmailAddress(emailAddress.String()), nil
}

func RebuildConfirmedEmailAddress(input string) ConfirmedEmailAddress {
	return ConfirmedEmailAddress(input)
}

func (emailAddress ConfirmedEmailAddress) String() string {
	return string(emailAddress)
}

func (emailAddress ConfirmedEmailAddress) Equals(other EmailAddress) bool {
	return emailAddress.String() == other.String()
}
