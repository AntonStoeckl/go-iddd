package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

func ConfirmEmailAddress(eventStream es.EventStream, command domain.ConfirmCustomerEmailAddress) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if err := assertNotDeleted(customer); err != nil {
		return nil, errors.Wrap(err, "confirmEmailAddress")
	}

	if err := assertMatchingConfirmationHash(customer.emailAddressConfirmationHash, command.ConfirmationHash()); err != nil {
		event := domain.BuildCustomerEmailAddressConfirmationFailed(
			command.CustomerID(),
			command.ConfirmationHash(),
			err,
			command.MessageID(),
			customer.currentStreamVersion+1,
		)

		return es.RecordedEvents{event}, nil
	}

	switch customer.emailAddress.(type) {
	case value.UnconfirmedEmailAddress:
		event := domain.BuildCustomerEmailAddressConfirmed(
			command.CustomerID(),
			value.ToConfirmedEmailAddress(customer.emailAddress),
			command.MessageID(),
			customer.currentStreamVersion+1,
		)

		return es.RecordedEvents{event}, nil
	case value.ConfirmedEmailAddress:
		return nil, nil
	default:
		// until Go has "union types" we need to use an interface and this case could exist - we don't want to hide it
		panic("ConfirmEmailAddress(): emailAddress is neither UnconfirmedEmailAddress nor ConfirmedEmailAddress")
	}
}
