package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

const failureReasonWrongHash = "wrong confirmation hash supplied"

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmCustomerEmailAddress) (es.DomainEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if !wasNotDeleted(customer) {
		return nil, errors.Wrap(wasDeletedErr, "confirmEmailAddress")
	}

	if !hasSuppliedMatchingConfirmationHash(customer.emailAddressConfirmationHash, command.ConfirmationHash()) {
		event := events.BuildCustomerEmailAddressConfirmationFailed(
			customer.id,
			customer.emailAddress,
			command.ConfirmationHash(),
			failureReasonWrongHash,
			customer.currentStreamVersion+1,
		)

		return es.DomainEvents{event}, nil
	}

	if customer.isEmailAddressConfirmed {
		return nil, nil
	}

	event := events.BuildCustomerEmailAddressConfirmed(
		customer.id,
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.DomainEvents{event}, nil
}
