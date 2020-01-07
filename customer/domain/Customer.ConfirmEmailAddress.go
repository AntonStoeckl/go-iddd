package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

func ConfirmEmailAddress(eventStream shared.DomainEvents, command *commands.ConfirmEmailAddress) shared.DomainEvents {
	var confirmableEmailAddress *values.ConfirmableEmailAddress
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case *events.Registered:
			confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
		case *events.EmailAddressConfirmed:
			confirmableEmailAddress = confirmableEmailAddress.MarkAsConfirmed()
		case *events.EmailAddressChanged:
			confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
		}

		currentStreamVersion = event.StreamVersion()
	}

	if confirmableEmailAddress.IsConfirmed() {
		return nil
	}

	if !confirmableEmailAddress.IsConfirmedBy(command.ConfirmationHash()) {
		event := events.EmailAddressConfirmationHasFailed(
			command.CustomerID(),
			command.ConfirmationHash(),
			currentStreamVersion+1,
		)

		return shared.DomainEvents{event}
	}

	event := events.EmailAddressWasConfirmed(
		command.CustomerID(),
		command.EmailAddress(),
		currentStreamVersion+1,
	)

	return shared.DomainEvents{event}
}
