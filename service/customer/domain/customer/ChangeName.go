package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func ChangeName(eventStream es.DomainEvents, command commands.ChangeCustomerName) es.DomainEvents {
	state := buildCustomerStateFrom(eventStream)

	if state.personName.Equals(command.PersonName()) {
		return nil
	}

	event := events.CustomerNameWasChanged(
		state.id,
		command.PersonName(),
		state.currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
