package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ChangeEmailAddress(eventStream es.EventStream, command commands.ChangeCustomerEmailAddress) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if !wasNotDeleted(customer) {
		return nil, errors.Wrap(wasDeletedErr, "changeEmailAddress")
	}

	if customer.emailAddress.Equals(command.EmailAddress()) {
		return nil, nil
	}

	event := events.BuildCustomerEmailAddressChanged(
		customer.id,
		command.EmailAddress(),
		command.ConfirmationHash(),
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}, nil
}
