package domain

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

type ConfirmEmailAddress struct {
	id               *valueobjects.ID
	emailAddress     *valueobjects.EmailAddress
	confirmationHash *valueobjects.ConfirmationHash
}

/*** Factory Method ***/

func NewConfirmEmailAddress(
	id *valueobjects.ID,
	emailAddress *valueobjects.EmailAddress,
	confirmationHash *valueobjects.ConfirmationHash,
) (*ConfirmEmailAddress, error) {

	command := &ConfirmEmailAddress{
		id:               id,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	if err := shared.AssertAllPropertiesAreNotNil(command); err != nil {
		return nil, err
	}

	return command, nil
}

/*** Getter Methods ***/

func (confirmEmailAddress *ConfirmEmailAddress) ID() *valueobjects.ID {
	return confirmEmailAddress.id
}

func (confirmEmailAddress *ConfirmEmailAddress) EmailAddress() *valueobjects.EmailAddress {
	return confirmEmailAddress.emailAddress
}

func (confirmEmailAddress *ConfirmEmailAddress) ConfirmationHash() *valueobjects.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

/*** Implement shared.Command ***/

func (confirmEmailAddress *ConfirmEmailAddress) AggregateIdentifier() shared.AggregateIdentifier {
	return confirmEmailAddress.id
}

func (confirmEmailAddress *ConfirmEmailAddress) CommandName() string {
	return shared.BuildCommandNameFor(confirmEmailAddress)
}
