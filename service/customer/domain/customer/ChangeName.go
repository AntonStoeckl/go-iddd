package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ChangeName(eventStream es.DomainEvents, command commands.ChangeCustomerName) (es.DomainEvents, error) {
	state := buildCustomerStateFrom(eventStream)

	if err := MustNotBeDeleted(state); err != nil {
		return nil, errors.Wrap(err, "changeCustomerName")
	}

	if state.personName.Equals(command.PersonName()) {
		return nil, nil
	}

	event := events.CustomerNameWasChanged(
		state.id,
		command.PersonName(),
		state.currentStreamVersion+1,
	)

	return es.DomainEvents{event}, nil
}
