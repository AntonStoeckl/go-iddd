package valueobjects

import (
	"errors"
)

type ConfirmableEmailAddress interface {
	Confirm(with ConfirmationHash) (*confirmableEmailAddress, error)
	String() string
	Equals(other EmailAddress) bool
	IsConfirmed() bool
}

type confirmableEmailAddress struct {
	value            string
	confirmationHash ConfirmationHash
	isConfirmed      bool
}

func NewConfirmableEmailAddress(from string) *confirmableEmailAddress {
	newEmailAddress := newConfirmableEmailAddress(from, GenerateConfirmationHash(from))
	// TODO: validation

	return newEmailAddress
}

func newConfirmableEmailAddress(from string, with ConfirmationHash) *confirmableEmailAddress {
	return &confirmableEmailAddress{
		value:            from,
		confirmationHash: with,
	}
}

func ReconstituteConfirmableEmailAddress(from string, withConfirmationHash string) *confirmableEmailAddress {
	return newConfirmableEmailAddress(from, ReconstituteConfirmationHash(withConfirmationHash))
}

func (confirmableEmailAddress *confirmableEmailAddress) Confirm(with ConfirmationHash) (*confirmableEmailAddress, error) {
	if confirmableEmailAddress.confirmationHash.Equals(with) {
		return nil, errors.New("confirmableEmailAddress - confirmationHash does not match")
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
