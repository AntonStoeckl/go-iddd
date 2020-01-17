package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

func ChangeEmailAddress(eventStream shared.DomainEvents, command commands.ChangeEmailAddress) shared.DomainEvents {
	var emailAddress values.EmailAddress
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.Registered:
			emailAddress = actualEvent.EmailAddress()
		case events.EmailAddressChanged:
			emailAddress = actualEvent.EmailAddress()
		}

		currentStreamVersion = event.StreamVersion()
	}

	if emailAddress.Equals(command.EmailAddress()) {
		return nil
	}

	event := events.EmailAddressWasChanged(
		command.CustomerID(),
		command.EmailAddress(),
		values.GenerateConfirmationHash(command.EmailAddress().EmailAddress()),
		currentStreamVersion+1,
	)

	return shared.DomainEvents{event}
}
