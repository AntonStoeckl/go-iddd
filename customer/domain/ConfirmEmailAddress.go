package domain

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

type ConfirmEmailAddress interface {
	ID() valueobjects.ID
	EmailAddress() valueobjects.EmailAddress
	ConfirmationHash() valueobjects.ConfirmationHash

	shared.Command
}

type confirmEmailAddress struct {
	id               valueobjects.ID
	emailAddress     valueobjects.EmailAddress
	confirmationHash valueobjects.ConfirmationHash
}

func NewConfirmEmailAddress(
	id valueobjects.ID,
	emailAddress valueobjects.EmailAddress,
	confirmationHash valueobjects.ConfirmationHash,
) (*confirmEmailAddress, error) {

	command := &confirmEmailAddress{
		id:               id,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	if err := shared.AssertAllPropertiesAreNotNil(command); err != nil {
		return nil, err
	}

	return command, nil
}

func (confirmEmailAddress *confirmEmailAddress) ID() valueobjects.ID {
	return confirmEmailAddress.id
}

func (confirmEmailAddress *confirmEmailAddress) EmailAddress() valueobjects.EmailAddress {
	return confirmEmailAddress.emailAddress
}

func (confirmEmailAddress *confirmEmailAddress) ConfirmationHash() valueobjects.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

func (confirmEmailAddress *confirmEmailAddress) Identifier() string {
	return confirmEmailAddress.id.String()
}

func (confirmEmailAddress *confirmEmailAddress) CommandName() string {
	return shared.BuildNameFor(confirmEmailAddress)
}
