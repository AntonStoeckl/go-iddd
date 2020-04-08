package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func Delete(eventStream es.DomainEvents) es.DomainEvents {
	customer := buildCurrentStateFrom(eventStream)

	if !wasNotDeleted(customer) {
		return nil
	}

	event := events.CustomerWasDeleted(
		customer.id,
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
