package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type state struct {
	id                           values.CustomerID
	personName                   values.PersonName
	emailAddress                 values.EmailAddress
	emailAddressConfirmationHash values.ConfirmationHash
	isEmailAddressConfirmed      bool
	isDeleted                    bool
	currentStreamVersion         uint
}

func buildCustomerStateFrom(eventStream es.DomainEvents) state {
	state := state{}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			state.id = actualEvent.CustomerID()
			state.personName = actualEvent.PersonName()
			state.emailAddress = actualEvent.EmailAddress()
			state.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
		case events.CustomerEmailAddressConfirmed:
			state.isEmailAddressConfirmed = true
		case events.CustomerEmailAddressChanged:
			state.emailAddress = actualEvent.EmailAddress()
			state.emailAddressConfirmationHash = actualEvent.ConfirmationHash()
			state.isEmailAddressConfirmed = false
		case events.CustomerNameChanged:
			state.personName = actualEvent.PersonName()
		case events.CustomerDeleted:
			state.isDeleted = true
		}

		state.currentStreamVersion = event.StreamVersion()
	}

	return state
}
