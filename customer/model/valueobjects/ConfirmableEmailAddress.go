package valueobjects

import (
	"errors"
	"regexp"
)

var (
	confirmableEmailAddressRegExp = regexp.MustCompile(`^[^\s]+@[^\s]+\.[\w]{2,}$`)
)

type ConfirmableEmailAddress interface {
	Confirm(given ConfirmationHash) (*confirmableEmailAddress, error)
	IsConfirmed() bool

	EmailAddress
}

type confirmableEmailAddress struct {
	value            string
	confirmationHash ConfirmationHash
	isConfirmed      bool
}

func NewConfirmableEmailAddress(from string) (*confirmableEmailAddress, error) {
	newEmailAddress := newConfirmableEmailAddress(from, GenerateConfirmationHash(from))

	if err := newEmailAddress.mustBeValid(); err != nil {
		return nil, err
	}

	return newEmailAddress, nil
}

func newConfirmableEmailAddress(from string, with ConfirmationHash) *confirmableEmailAddress {
	return &confirmableEmailAddress{
		value:            from,
		confirmationHash: with,
	}
}

func (confirmableEmailAddress *confirmableEmailAddress) mustBeValid() error {
	if matched := confirmableEmailAddressRegExp.MatchString(confirmableEmailAddress.value); matched != true {
		return errors.New("confirmableEmailAddress - invalid input given")
	}

	return nil
}

func ReconstituteConfirmableEmailAddress(from string, withConfirmationHash string) *confirmableEmailAddress {
	return newConfirmableEmailAddress(from, ReconstituteConfirmationHash(withConfirmationHash))
}

func (confirmableEmailAddress *confirmableEmailAddress) Confirm(given ConfirmationHash) (*confirmableEmailAddress, error) {
	if err := confirmableEmailAddress.confirmationHash.MustMatch(given); err != nil {
		return nil, err
	}

	confirmedEmailAddress := newConfirmableEmailAddress(
		confirmableEmailAddress.value,
		confirmableEmailAddress.confirmationHash,
	)

	confirmedEmailAddress.isConfirmed = true

	return confirmedEmailAddress, nil
}

func (confirmableEmailAddress *confirmableEmailAddress) String() string {
	return confirmableEmailAddress.value
}

func (confirmableEmailAddress *confirmableEmailAddress) Equals(other EmailAddress) bool {
	return confirmableEmailAddress.String() == other.String()
}

func (confirmableEmailAddress *confirmableEmailAddress) IsConfirmed() bool {
	return confirmableEmailAddress.isConfirmed
}
