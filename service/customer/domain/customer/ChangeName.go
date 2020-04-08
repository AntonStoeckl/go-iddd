package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ChangeName(eventStream es.DomainEvents, command commands.ChangeCustomerName) (es.DomainEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if !wasNotDeleted(customer) {
		return nil, errors.Wrap(wasDeletedErr, "changeCustomerName")
	}

	if customer.personName.Equals(command.PersonName()) {
		return nil, nil
	}

	event := events.CustomerNameWasChanged(
		customer.id,
		command.PersonName(),
		customer.currentStreamVersion+1,
	)

	return es.DomainEvents{event}, nil
}
