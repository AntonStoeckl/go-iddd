package domain

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
)

func ChangeCustomerEmailAddress(eventStream es.DomainEvents, command commands.ChangeEmailAddress) es.DomainEvents {
	var emailAddress values.EmailAddress
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			emailAddress = actualEvent.EmailAddress()
		case events.CustomerEmailAddressChanged:
			emailAddress = actualEvent.EmailAddress()
		}

		currentStreamVersion = event.StreamVersion()
	}

	if emailAddress.Equals(command.EmailAddress()) {
		return nil
	}

	event := events.CustomerEmailAddressWasChanged(
		command.CustomerID(),
		command.EmailAddress(),
		command.ConfirmationHash(),
		currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
