package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

func ChangeEmailAddress(eventStream es.EventStream, command domain.ChangeCustomerEmailAddress) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if err := assertNotDeleted(customer); err != nil {
		return nil, errors.Wrap(err, "changeEmailAddress")
	}

	if customer.emailAddress.Equals(command.EmailAddress()) {
		return nil, nil
	}

	event := domain.BuildCustomerEmailAddressChanged(
		customer.id,
		command.EmailAddress(),
		command.ConfirmationHash(),
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}, nil
}
