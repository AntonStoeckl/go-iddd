package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func ChangeEmailAddress(eventStream es.DomainEvents, command commands.ChangeCustomerEmailAddress) es.DomainEvents {
	state := buildCustomerStateFrom(eventStream)

	if state.emailAddress.Equals(command.EmailAddress()) {
		return nil
	}

	event := events.CustomerEmailAddressWasChanged(
		state.id,
		command.EmailAddress(),
		command.ConfirmationHash(),
		state.currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
