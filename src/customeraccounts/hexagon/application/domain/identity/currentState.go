package identity

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type currentState struct {
	id                   value.IdentityID
	emailAddress         value.EmailAddress
	password             value.HashedPassword
	isDeleted            bool
	currentStreamVersion uint
}

func buildCurrentStateFrom(eventStream es.EventStream) currentState {
	customer := currentState{}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case domain.IdentityRegistered:
			customer.id = actualEvent.IdentityID()
			customer.emailAddress = actualEvent.EmailAddress()
			customer.password = actualEvent.Password()
		default:
			// until Go has "sum types" we need to use an interface (Event) and this case could exist - we don't want to hide it
			panic("buildCurrentStateFrom(eventStream): unknown event " + event.Meta().EventName())
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	return customer
}
