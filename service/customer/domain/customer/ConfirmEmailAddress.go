package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmCustomerEmailAddress) es.DomainEvents {
	state := buildCustomerStateFrom(eventStream)

	if !state.emailAddressConfirmationHash.Equals(command.ConfirmationHash()) {
		event := events.CustomerEmailAddressConfirmationHasFailed(
			state.id,
			state.emailAddress,
			command.ConfirmationHash(),
			state.currentStreamVersion+1,
		)

		return es.DomainEvents{event}
	}

	if state.isEmailAddressConfirmed {
		return nil
	}

	event := events.CustomerEmailAddressWasConfirmed(
		state.id,
		state.emailAddress,
		state.currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
