package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

type ConfirmEmailAddress struct {
	id               *values.ID
	emailAddress     *values.EmailAddress
	confirmationHash *values.ConfirmationHash
}

/*** Factory Method ***/

func NewConfirmEmailAddress(
	id *values.ID,
	emailAddress *values.EmailAddress,
	confirmationHash *values.ConfirmationHash,
) (*ConfirmEmailAddress, error) {

	command := &ConfirmEmailAddress{
		id:               id,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	if err := shared.AssertAllCommandPropertiesAreNotNil(command); err != nil {
		return nil, xerrors.Errorf("confirmEmailAddress.New -> %s: %w", err, shared.ErrNilInput)
	}

	return command, nil
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
