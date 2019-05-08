package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

const emailAddressConfirmedAggregateName = "Customer"

type EmailAddressConfirmed struct {
	id           *valueobjects.ID
	emailAddress *valueobjects.EmailAddress

	meta *shared.DomainEventMeta
}

/*** Factory Methods ***/

func EmailAddressWasConfirmed(
	id *valueobjects.ID,
	emailAddress *valueobjects.EmailAddress,
) *EmailAddressConfirmed {

	emailAddressConfirmed := &EmailAddressConfirmed{
		id:           id,
		emailAddress: emailAddress,
	}

	emailAddressConfirmed.meta = shared.NewDomainEventMeta(
		id.String(),
		emailAddressConfirmed,
		emailAddressConfirmedAggregateName,
	)

	return emailAddressConfirmed
}

/*** Getter Methods ***/

func (emailAddressConfirmed *EmailAddressConfirmed) ID() *valueobjects.ID {
	return emailAddressConfirmed.id
}

func (emailAddressConfirmed *EmailAddressConfirmed) EmailAddress() *valueobjects.EmailAddress {
	return emailAddressConfirmed.emailAddress
}

/*** Implement shared.DomainEvent ***/

func (emailAddressConfirmed *EmailAddressConfirmed) Identifier() string {
	return emailAddressConfirmed.meta.Identifier
}

func (emailAddressConfirmed *EmailAddressConfirmed) EventName() string {
	return emailAddressConfirmed.meta.EventName
}

func (emailAddressConfirmed *EmailAddressConfirmed) OccurredAt() string {
	return emailAddressConfirmed.meta.OccurredAt
}

/*** Implement json.Marshaler ***/

func (emailAddressConfirmed *EmailAddressConfirmed) MarshalJSON() ([]byte, error) {
	data := &struct {
		ID           *valueobjects.ID           `json:"id"`
		EmailAddress *valueobjects.EmailAddress `json:"emailAddress"`
		Meta         *shared.DomainEventMeta    `json:"meta"`
	}{
		ID:           emailAddressConfirmed.id,
		EmailAddress: emailAddressConfirmed.emailAddress,
		Meta:         emailAddressConfirmed.meta,
	}

	return json.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (emailAddressConfirmed *EmailAddressConfirmed) UnmarshalJSON(data []byte) error {
	values := &struct {
		ID           *valueobjects.ID           `json:"id"`
		EmailAddress *valueobjects.EmailAddress `json:"emailAddress"`
		Meta         *shared.DomainEventMeta    `json:"meta"`
	}{}

	if err := json.Unmarshal(data, values); err != nil {
		return err
	}

	emailAddressConfirmed.id = values.ID
	emailAddressConfirmed.emailAddress = values.EmailAddress
	emailAddressConfirmed.meta = values.Meta

	return nil
}
