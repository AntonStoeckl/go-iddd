package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
	"github.com/cockroachdb/errors"
)

func ConfirmEmailAddress(eventStream es.EventStream, command domain.ConfirmCustomerEmailAddress) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if err := assertNotDeleted(customer); err != nil {
		return nil, errors.Wrap(err, "confirmEmailAddress")
	}

	if err := assertMatchingConfirmationHash(customer.emailAddressConfirmationHash, command.ConfirmationHash()); err != nil {
		event := domain.BuildCustomerEmailAddressConfirmationFailed(
			customer.id,
			customer.emailAddress,
			command.ConfirmationHash(),
			err,
			customer.currentStreamVersion+1,
		)

		return es.RecordedEvents{event}, nil
	}

	if customer.isEmailAddressConfirmed {
		return nil, nil
	}

	event := domain.BuildCustomerEmailAddressConfirmed(
		customer.id,
		customer.emailAddress,
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}, nil
}
