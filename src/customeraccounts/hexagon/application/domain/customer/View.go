package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
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

func BuildViewFrom(eventStream es.EventStream) View {
	customer := buildCurrentStateFrom(eventStream)

	customerView := View{
		ID:           customer.id.String(),
		EmailAddress: customer.emailAddress.String(),
		GivenName:    customer.personName.GivenName(),
		FamilyName:   customer.personName.FamilyName(),
		IsDeleted:    customer.isDeleted,
		Version:      customer.currentStreamVersion,
	}

	switch customer.emailAddress.(type) {
	case value.ConfirmedEmailAddress:
		customerView.IsEmailAddressConfirmed = true
	case value.UnconfirmedEmailAddress:
		customerView.IsEmailAddressConfirmed = false
	default:
		// until Go has "union types" we need to use an interface and this case could exist - we don't want to hide it
		panic("BuildViewFrom(eventStream): emailAddress is neither UnconfirmedEmailAddress nor ConfirmedEmailAddress")
	}

	return customerView
}
