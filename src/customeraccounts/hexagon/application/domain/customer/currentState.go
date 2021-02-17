package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type currentState struct {
	id                           value.CustomerID
	personName                   value.PersonName
	emailAddress                 value.EmailAddress
	emailAddressConfirmationHash value.ConfirmationHash
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
			customer.emailAddress = actualEvent.EmailAddress()
		case domain.CustomerEmailAddressChanged:
			customer.emailAddress = actualEvent.EmailAddress()
			customer.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
		case domain.CustomerNameChanged:
			customer.personName = actualEvent.PersonName()
		case domain.CustomerDeleted:
			customer.isDeleted = true
		case domain.CustomerEmailAddressConfirmationFailed:
			// nothing to project here
		default:
			// until Go has "sum types" we need to use an interface (Event) and this case could exist - we don't want to hide it
			panic("buildCurrentStateFrom(eventStream): unknown event " + event.Meta().EventName())
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	return customer
}
