package customer

import (
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

func BuildViewFrom(eventStream es.EventStream) View {
	customer := buildCurrentStateFrom(eventStream)

	customerView := View{
		ID:                      customer.id.String(),
		EmailAddress:            customer.emailAddress.String(),
		IsEmailAddressConfirmed: customer.isEmailAddressConfirmed,
		GivenName:               customer.personName.GivenName(),
		FamilyName:              customer.personName.FamilyName(),
		IsDeleted:               customer.isDeleted,
		Version:                 customer.currentStreamVersion,
	}

	return customerView
}
