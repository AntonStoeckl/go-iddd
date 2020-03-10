package customer

import (
	"go-iddd/service/customer/application/domain/customer/commands"
	"go-iddd/service/customer/application/domain/customer/events"
	"go-iddd/service/customer/application/domain/customer/values"
	"go-iddd/service/lib/es"
)

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmCustomerEmailAddress) es.DomainEvents {
	var emailAddress values.EmailAddress
	var confirmationHash values.ConfirmationHash
	var isConfirmed bool
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			emailAddress = actualEvent.EmailAddress()
			confirmationHash = actualEvent.ConfirmationHash()
		case events.CustomerEmailAddressConfirmed:
			isConfirmed = true
		case events.CustomerEmailAddressChanged:
			emailAddress = actualEvent.EmailAddress()
			confirmationHash = actualEvent.ConfirmationHash()
			isConfirmed = false
		}

		currentStreamVersion = event.StreamVersion()
	}

	if !confirmationHash.Equals(command.ConfirmationHash()) {
		event := events.CustomerEmailAddressConfirmationHasFailed(
			command.CustomerID(),
			emailAddress,
			command.ConfirmationHash(),
			currentStreamVersion+1,
		)

		return es.DomainEvents{event}
	}

	if isConfirmed {
		return nil
	}

	event := events.CustomerEmailAddressWasConfirmed(
		command.CustomerID(),
		emailAddress,
		currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
