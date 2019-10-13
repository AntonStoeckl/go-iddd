package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
)

func (customer *Customer) ConfirmEmailAddress(with *commands.ConfirmEmailAddress) error {
	if customer.confirmableEmailAddress.IsConfirmed() {
		return nil
	}

	err := customer.confirmableEmailAddress.ShouldConfirm(
		with.EmailAddress(),
		with.ConfirmationHash(),
	)

	if err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrDomainConstraintsViolation), "customer.ConfirmEmailAddress")
	}

	customer.recordThat(
		events.EmailAddressWasConfirmed(
			with.ID(),
			with.EmailAddress(),
			customer.currentStreamVersion+1,
		),
	)

	return nil
}

func (customer *Customer) whenEmailAddressWasConfirmed(actualEvent *events.EmailAddressConfirmed) {
	customer.confirmableEmailAddress = customer.confirmableEmailAddress.MarkAsConfirmed()
}
