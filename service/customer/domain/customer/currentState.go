package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type currentState struct {
	id                           values.CustomerID
	personName                   values.PersonName
	emailAddress                 values.EmailAddress
	emailAddressConfirmationHash values.ConfirmationHash
	isEmailAddressConfirmed      bool
	isDeleted                    bool
	currentStreamVersion         uint
}

func buildCurrentStateFrom(eventStream es.DomainEvents) currentState {
	customer := currentState{}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			customer.id = actualEvent.CustomerID()
			customer.personName = actualEvent.PersonName()
			customer.emailAddress = actualEvent.EmailAddress()
			customer.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
		case events.CustomerEmailAddressConfirmed:
			customer.isEmailAddressConfirmed = true
		case events.CustomerEmailAddressChanged:
			customer.emailAddress = actualEvent.EmailAddress()
			customer.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
			customer.isEmailAddressConfirmed = false
		case events.CustomerNameChanged:
			customer.personName = actualEvent.PersonName()
		case events.CustomerDeleted:
			customer.isDeleted = true
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	return customer
}
