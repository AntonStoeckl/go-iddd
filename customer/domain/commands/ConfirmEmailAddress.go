package commands

import (
	"errors"
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

	confirmEmailAddress := &ConfirmEmailAddress{
		id:               id,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	if err := confirmEmailAddress.shouldHaveNonEmpty(id, emailAddress, confirmationHash); err != nil {
		return nil, xerrors.Errorf("confirmEmailAddress.New: %s: %w", err, shared.ErrInputIsInvalid)
	}

	return confirmEmailAddress, nil
}

func (confirmEmailAddress *ConfirmEmailAddress) shouldHaveNonEmpty(
	id *values.ID,
	emailAddress *values.EmailAddress,
	confirmationHash *values.ConfirmationHash,
) error {

	if id == nil {
		return errors.New("id is nil")
	}

	if id.Equals(&values.ID{}) {
		return errors.New("id is empty (not created with factory method)")
	}

	if emailAddress == nil {
		return errors.New("emailAddress is nil")
	}

	if emailAddress.Equals(&values.EmailAddress{}) {
		return errors.New("emailAddress is empty (not created with factory method)")
	}

	if confirmationHash == nil {
		return errors.New("confirmationHash is nil")
	}

	if confirmationHash.Equals(&values.ConfirmationHash{}) {
		return errors.New("confirmationHash is empty (not created with factory method)")
	}

	return nil
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
