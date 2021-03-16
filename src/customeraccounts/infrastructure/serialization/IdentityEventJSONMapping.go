package serialization

import "github.com/AntonStoeckl/go-iddd/src/shared/es"

type IdentityRegisteredForJSON struct {
	IdentityID       string              `json:"identityID"`
	EmailAddress     string              `json:"emailAddress"`
	ConfirmationHash string              `json:"confirmationHash"`
	Password         string              `json:"password"`
	Meta             es.EventMetaForJSON `json:"meta"`
}

type IdentityDeletedForJSON struct {
	IdentityID string              `json:"identityID"`
	Meta       es.EventMetaForJSON `json:"meta"`
}
