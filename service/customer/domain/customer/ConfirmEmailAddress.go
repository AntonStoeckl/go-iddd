package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

const failureReasonWrongHash = "wrong confirmation hash supplied"

func ConfirmEmailAddress(eventStream es.DomainEvents, command commands.ConfirmCustomerEmailAddress) (es.DomainEvents, error) {
	state := buildCustomerStateFrom(eventStream)

	if err := MustNotBeDeleted(state); err != nil {
		return nil, errors.Wrap(err, "confirmEmailAddress")
	}

	if !IsMatchingConfirmationHash(state.emailAddressConfirmationHash, command.ConfirmationHash()) {
		event := events.CustomerEmailAddressConfirmationHasFailed(
			state.id,
			state.emailAddress,
			command.ConfirmationHash(),
			failureReasonWrongHash,
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
