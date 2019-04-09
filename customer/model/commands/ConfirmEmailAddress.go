package commands

import (
	"go-iddd/customer/model/valueobjects"
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

func (register *confirmEmailAddress) ID() valueobjects.ID {
	return register.id
}

func (register *confirmEmailAddress) EmailAddress() valueobjects.EmailAddress {
	return register.emailAddress
}

func (register *confirmEmailAddress) ConfirmationHash() valueobjects.ConfirmationHash {
	return register.confirmationHash
}

func (register *confirmEmailAddress) Identifier() string {
	return register.id.String()
}

func (register *confirmEmailAddress) CommandName() string {
	return shared.BuildNameFor(register)
}
