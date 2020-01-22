package domain

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
)

func ChangeEmailAddress(eventStream lib.DomainEvents, command commands.ChangeEmailAddress) lib.DomainEvents {
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
		command.ConfirmationHash(),
		currentStreamVersion+1,
	)

	return lib.DomainEvents{event}
}
