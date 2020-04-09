package serialization

import "github.com/AntonStoeckl/go-iddd/service/lib/es"

type CustomerRegisteredForJSON struct {
	CustomerID       string       `json:"customerID"`
	EmailAddress     string       `json:"emailAddress"`
	ConfirmationHash string       `json:"confirmationHash"`
	PersonGivenName  string       `json:"personGivenName"`
	PersonFamilyName string       `json:"personFamilyName"`
	Meta             es.EventMeta `json:"meta"`
}

type CustomerEmailAddressConfirmedForJSON struct {
	CustomerID   string       `json:"customerID"`
	EmailAddress string       `json:"emailAddress"`
	Meta         es.EventMeta `json:"meta"`
}

type CustomerEmailAddressConfirmationFailedForJSON struct {
	CustomerID       string       `json:"customerID"`
	EmailAddress     string       `json:"emailAddress"`
	ConfirmationHash string       `json:"confirmationHash"`
	Reason           string       `json:"reason"`
	Meta             es.EventMeta `json:"meta"`
}

type CustomerEmailAddressChangedForJSON struct {
	CustomerID           string       `json:"customerID"`
	EmailAddress         string       `json:"emailAddress"`
	ConfirmationHash     string       `json:"confirmationHash"`
	PreviousEmailAddress string       `json:"previousEmailAddress"`
	Meta                 es.EventMeta `json:"meta"`
}

type CustomerNameChangedForJSON struct {
	CustomerID string       `json:"customerID"`
	GivenName  string       `json:"givenName"`
	FamilyName string       `json:"familyName"`
	Meta       es.EventMeta `json:"meta"`
}

type CustomerDeletedForJSON struct {
	CustomerID   string       `json:"customerID"`
	EmailAddress string       `json:"emailAddress"`
	Meta         es.EventMeta `json:"meta"`
}
