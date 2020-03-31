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

func BuildViewFrom(eventStream es.DomainEvents) View {
	state := buildCustomerStateFrom(eventStream)

	customerView := View{
		ID:                      state.id.String(),
		EmailAddress:            state.emailAddress.String(),
		IsEmailAddressConfirmed: state.isEmailAddressConfirmed,
		GivenName:               state.personName.GivenName(),
		FamilyName:              state.personName.FamilyName(),
		IsDeleted:               state.isDeleted,
		Version:                 state.currentStreamVersion,
	}

	return customerView
}
