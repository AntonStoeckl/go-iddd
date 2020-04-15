package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ConfirmEmailAddress(eventStream es.EventStream, command commands.ConfirmCustomerEmailAddress) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if err := assertNotDeleted(customer); err != nil {
		return nil, errors.Wrap(err, "confirmEmailAddress")
	}

	if err := assertMatchingConfirmationHash(customer.emailAddressConfirmationHash, command.ConfirmationHash()); err != nil {
		event := events.BuildCustomerEmailAddressConfirmationFailed(
			customer.id,
			customer.emailAddress,
			command.ConfirmationHash(),
			err,
			customer.currentStreamVersion+1,
		)

		return es.RecordedEvents{event}, nil
	}

	if customer.isEmailAddressConfirmed {
		return nil, nil
	}

	event := events.BuildCustomerEmailAddressConfirmed(
		customer.id,
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}, nil
}
