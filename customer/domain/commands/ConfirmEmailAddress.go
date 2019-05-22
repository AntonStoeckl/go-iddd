package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type ConfirmEmailAddress struct {
	id               *values.ID
	emailAddress     *values.EmailAddress
	confirmationHash *values.ConfirmationHash
}

/*** Factory Method ***/

func NewConfirmEmailAddress(
	id string,
	emailAddress string,
	confirmationHash string,
) (*ConfirmEmailAddress, error) {

	idValue, err := values.NewID(id)
	if err != nil {
		return nil, err
	}

	emailAddressValue, err := values.NewEmailAddress(emailAddress)
	if err != nil {
		return nil, err
	}

	confirmationHashValue, err := values.NewConfirmationHash(confirmationHash)
	if err != nil {
		return nil, err
	}

	confirmEmailAddress := &ConfirmEmailAddress{
		id:               idValue,
		emailAddress:     emailAddressValue,
		confirmationHash: confirmationHashValue,
	}

	return confirmEmailAddress, nil
}

/*** Getter Methods ***/

func (confirmEmailAddress *ConfirmEmailAddress) ID() *values.ID {
	return confirmEmailAddress.id
}

func (confirmEmailAddress *ConfirmEmailAddress) EmailAddress() *values.EmailAddress {
	return confirmEmailAddress.emailAddress
}

func (confirmEmailAddress *ConfirmEmailAddress) ConfirmationHash() *values.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

/*** Implement shared.Command ***/

func (confirmEmailAddress *ConfirmEmailAddress) AggregateIdentifier() shared.AggregateIdentifier {
	return confirmEmailAddress.id
}

func (confirmEmailAddress *ConfirmEmailAddress) CommandName() string {
	return shared.BuildCommandNameFor(confirmEmailAddress)
}
