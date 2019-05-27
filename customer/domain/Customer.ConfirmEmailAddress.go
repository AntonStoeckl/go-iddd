package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

func (customer *customer) confirmEmailAddress(with *commands.ConfirmEmailAddress) error {
	if customer.confirmableEmailAddress.IsConfirmed() {
		return nil
	}

	err := customer.confirmableEmailAddress.ShouldConfirm(
		with.EmailAddress(),
		with.ConfirmationHash(),
	)

	if err != nil {
		return xerrors.Errorf("customer.confirmEmailAddress -> %s: %w", err, shared.ErrDomainConstraintsViolation)
	}

	customer.recordThat(
		events.EmailAddressWasConfirmed(
			with.ID(),
			with.EmailAddress(),
		),
	)

	return nil
}

func (customer *customer) whenEmailAddressWasConfirmed(actualEvent *events.EmailAddressConfirmed) {
	customer.confirmableEmailAddress = customer.confirmableEmailAddress.MarkAsConfirmed()
}
