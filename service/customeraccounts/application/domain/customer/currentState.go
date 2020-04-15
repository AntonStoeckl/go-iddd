package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type currentState struct {
	id                           value.CustomerID
	personName                   value.PersonName
	emailAddress                 value.EmailAddress
	emailAddressConfirmationHash value.ConfirmationHash
	isEmailAddressConfirmed      bool
	isDeleted                    bool
	currentStreamVersion         uint
}

func buildCurrentStateFrom(eventStream es.EventStream) currentState {
	customer := currentState{}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case domain.CustomerRegistered:
			customer.id = actualEvent.CustomerID()
			customer.personName = actualEvent.PersonName()
			customer.emailAddress = actualEvent.EmailAddress()
			customer.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
		case domain.CustomerEmailAddressConfirmed:
			customer.isEmailAddressConfirmed = true
		case domain.CustomerEmailAddressChanged:
			customer.emailAddress = actualEvent.EmailAddress()
			customer.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
			customer.isEmailAddressConfirmed = false
		case domain.CustomerNameChanged:
			customer.personName = actualEvent.PersonName()
		case domain.CustomerDeleted:
			customer.isDeleted = true
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	return customer
}
