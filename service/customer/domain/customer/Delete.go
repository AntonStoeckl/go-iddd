package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func Delete(eventStream es.DomainEvents, command commands.DeleteCustomer) es.DomainEvents {
	var isDeleted bool
	var currentStreamVersion uint

	for _, event := range eventStream {
		switch event.(type) {
		case events.CustomerDeleted:
			isDeleted = true
		}

		currentStreamVersion = event.StreamVersion()
	}

	if isDeleted {
		return nil
	}

	event := events.CustomerWasDeleted(
		command.CustomerID(),
		currentStreamVersion+1,
	)

	return es.DomainEvents{event}
}
