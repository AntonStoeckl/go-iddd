package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ChangeEmailAddress(eventStream es.DomainEvents, command commands.ChangeCustomerEmailAddress) (es.DomainEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if !wasNotDeleted(customer) {
		return nil, errors.Wrap(wasDeletedErr, "changeEmailAddress")
	}

	if customer.emailAddress.Equals(command.EmailAddress()) {
		return nil, nil
	}

	event := events.CustomerEmailAddressWasChanged(
		customer.id,
		command.EmailAddress(),
		command.ConfirmationHash(),
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.DomainEvents{event}, nil
}
