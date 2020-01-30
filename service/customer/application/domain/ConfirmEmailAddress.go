package domain

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
)

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmEmailAddress) es.DomainEvents {
	var emailAddress values.EmailAddress
	var confirmationHash values.ConfirmationHash
	var isConfirmed bool
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.Registered:
			emailAddress = actualEvent.EmailAddress()
			confirmationHash = actualEvent.ConfirmationHash()
		case events.EmailAddressConfirmed:
			isConfirmed = true
		case events.EmailAddressChanged:
			emailAddress = actualEvent.EmailAddress()
			confirmationHash = actualEvent.ConfirmationHash()
			isConfirmed = false
		}

		currentStreamVersion = event.StreamVersion()
	}

	if isConfirmed {
		return nil
	}

	if !confirmationHash.Equals(command.ConfirmationHash()) {
		event := events.EmailAddressConfirmationHasFailed(
			command.CustomerID(),
			emailAddress,
			command.ConfirmationHash(),
			currentStreamVersion+1,
		)

		return es.DomainEvents{event}
	}

	event := events.EmailAddressWasConfirmed(
		command.CustomerID(),
		emailAddress,
		currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
