package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type View struct {
	ID                      string
	EmailAddress            string
	IsEmailAddressConfirmed bool
	GivenName               string
	FamilyName              string
	IsDeleted               bool
	Version                 uint
}

func BuildViewFrom(eventStream es.DomainEvents) View {
	customerView := View{}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			customerView.ID = actualEvent.CustomerID().ID()
			customerView.EmailAddress = actualEvent.EmailAddress().EmailAddress()
			customerView.GivenName = actualEvent.PersonName().GivenName()
			customerView.FamilyName = actualEvent.PersonName().FamilyName()
		case events.CustomerEmailAddressConfirmed:
			customerView.IsEmailAddressConfirmed = true
		case events.CustomerEmailAddressChanged:
			customerView.EmailAddress = actualEvent.EmailAddress().EmailAddress()
			customerView.IsEmailAddressConfirmed = false
		case events.CustomerNameChanged:
			customerView.GivenName = actualEvent.PersonName().GivenName()
			customerView.FamilyName = actualEvent.PersonName().FamilyName()
		case events.CustomerDeleted:
			customerView.IsDeleted = true
		}

		customerView.Version = event.StreamVersion()
	}

	return customerView
}
