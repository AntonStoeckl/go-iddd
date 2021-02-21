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

	switch actualEmailAddress := customer.emailAddress.(type) {
	case value.ConfirmedEmailAddress:
		return nil, nil
	case value.UnconfirmedEmailAddress:
		err := assertMatchingConfirmationHash(actualEmailAddress.ConfirmationHash(), command.ConfirmationHash())
		if err != nil {
			return es.RecordedEvents{
				domain.BuildCustomerEmailAddressConfirmationFailed(
					command.CustomerID(),
					command.ConfirmationHash(),
					err,
					command.MessageID(),
					customer.currentStreamVersion+1,
				),
			}, nil
		}

		return es.RecordedEvents{
			domain.BuildCustomerEmailAddressConfirmed(
				command.CustomerID(),
				value.ToConfirmedEmailAddress(customer.emailAddress),
				command.MessageID(),
				customer.currentStreamVersion+1,
			),
		}, nil
	default:
		// until Go has "union types" we need to use an interface and this case could exist - we don't want to hide it
		panic("ConfirmEmailAddress(): emailAddress is neither UnconfirmedEmailAddress nor ConfirmedEmailAddress")
	}
}
