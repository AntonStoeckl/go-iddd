package domain

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
)

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmEmailAddress) es.DomainEvents {
	var confirmationHash values.ConfirmationHash
	var isConfirmed bool
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			confirmationHash = actualEvent.ConfirmationHash()
		case events.CustomerEmailAddressConfirmed:
			isConfirmed = true
		case events.CustomerEmailAddressChanged:
			confirmationHash = actualEvent.ConfirmationHash()
			isConfirmed = false
		}

		currentStreamVersion = event.StreamVersion()
	}

	if !confirmationHash.Equals(command.ConfirmationHash()) {
		event := events.CustomerEmailAddressConfirmationHasFailed(
			command.CustomerID(),
			command.EmailAddress(),
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
		command.EmailAddress(),
		currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
