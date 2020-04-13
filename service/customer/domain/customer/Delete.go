package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func Delete(eventStream es.EventStream) es.RecordedEvents {
	customer := buildCurrentStateFrom(eventStream)

	if !wasNotDeleted(customer) {
		return nil
	}

	event := events.BuildCustomerDeleted(
		customer.id,
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}
}
