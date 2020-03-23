package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func Delete(eventStream es.DomainEvents) es.DomainEvents {
	state := buildCustomerStateFrom(eventStream)

	if state.isDeleted {
		return nil
	}

	event := events.CustomerWasDeleted(
		state.id,
		state.currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
