package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

func ConfirmEmailAddress(eventStream shared.DomainEvents, command commands.ConfirmEmailAddress) shared.DomainEvents {
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

		return shared.DomainEvents{event}
	}

	event := events.EmailAddressWasConfirmed(
		command.CustomerID(),
		emailAddress,
		currentStreamVersion+1,
	)

	return shared.DomainEvents{event}
}
