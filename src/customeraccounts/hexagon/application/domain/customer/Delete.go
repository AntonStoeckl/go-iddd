package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

func Delete(eventStream es.EventStream, command domain.DeleteCustomer) es.RecordedEvents {
	customer := buildCurrentStateFrom(eventStream)

	if err := assertNotDeleted(customer); err != nil {
		return nil
	}

	event := domain.BuildCustomerDeleted(
		command.CustomerID(),
		customer.emailAddress,
		command.MessageID(),
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}
}
