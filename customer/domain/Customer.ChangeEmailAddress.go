package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

func ChangeEmailAddress(eventStream shared.DomainEvents, command *commands.ChangeEmailAddress) shared.DomainEvents {
	var confirmableEmailAddress *values.ConfirmableEmailAddress
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case *events.Registered:
			confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
		case *events.EmailAddressChanged:
			confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
		}

		currentStreamVersion = event.StreamVersion()
	}

	if confirmableEmailAddress.Equals(command.EmailAddress()) {
		return nil
	}

	event := events.EmailAddressWasChanged(
		command.CustomerID(),
		command.EmailAddress().ToConfirmable(),
		currentStreamVersion+1,
	)

	return shared.DomainEvents{event}
}
