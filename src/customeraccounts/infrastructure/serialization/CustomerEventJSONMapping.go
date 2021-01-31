package serialization

import "github.com/AntonStoeckl/go-iddd/src/shared/es"

type CustomerRegisteredForJSON struct {
	CustomerID       string              `json:"customerID"`
	EmailAddress     string              `json:"emailAddress"`
	ConfirmationHash string              `json:"confirmationHash"`
	PersonGivenName  string              `json:"personGivenName"`
	PersonFamilyName string              `json:"personFamilyName"`
	Meta             es.EventMetaForJSON `json:"meta"`
}

type CustomerEmailAddressConfirmedForJSON struct {
	CustomerID   string              `json:"customerID"`
	EmailAddress string              `json:"emailAddress"`
	Meta         es.EventMetaForJSON `json:"meta"`
}

type CustomerEmailAddressConfirmationFailedForJSON struct {
	CustomerID       string              `json:"customerID"`
	EmailAddress     string              `json:"emailAddress"`
	ConfirmationHash string              `json:"confirmationHash"`
	Reason           string              `json:"reason"`
	Meta             es.EventMetaForJSON `json:"meta"`
}

type CustomerEmailAddressChangedForJSON struct {
	CustomerID           string              `json:"customerID"`
	EmailAddress         string              `json:"emailAddress"`
	ConfirmationHash     string              `json:"confirmationHash"`
	PreviousEmailAddress string              `json:"previousEmailAddress"`
	Meta                 es.EventMetaForJSON `json:"meta"`
}

type CustomerNameChangedForJSON struct {
	CustomerID string              `json:"customerID"`
	GivenName  string              `json:"givenName"`
	FamilyName string              `json:"familyName"`
	Meta       es.EventMetaForJSON `json:"meta"`
}

type CustomerDeletedForJSON struct {
	CustomerID   string              `json:"customerID"`
	EmailAddress string              `json:"emailAddress"`
	Meta         es.EventMetaForJSON `json:"meta"`
}
