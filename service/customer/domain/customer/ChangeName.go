package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ChangeName(eventStream es.DomainEvents, command commands.ChangeCustomerName) (es.DomainEvents, error) {
	state := buildCustomerStateFrom(eventStream)

	if state.isDeleted {
		err := errors.New("customer is deleted")

		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, "changeCustomerName")
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
