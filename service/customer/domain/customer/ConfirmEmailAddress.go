package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmCustomerEmailAddress) (es.DomainEvents, error) {
	state := buildCustomerStateFrom(eventStream)

	if state.isDeleted {
		err := errors.New("customer is deleted")

		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, "confirmEmailAddress")
	}

	if !state.emailAddressConfirmationHash.Equals(command.ConfirmationHash()) {
		event := events.CustomerEmailAddressConfirmationHasFailed(
			state.id,
			state.emailAddress,
			command.ConfirmationHash(),
			state.currentStreamVersion+1,
		)

		return es.DomainEvents{event}, nil
	}

	if state.isEmailAddressConfirmed {
		return nil, nil
	}

	event := events.CustomerEmailAddressWasConfirmed(
		state.id,
		state.emailAddress,
		state.currentStreamVersion+1,
	)

	return es.DomainEvents{event}, nil
}
