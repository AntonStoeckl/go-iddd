package events

import (
	"encoding/json"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

const emailAddressConfirmedAggregateName = "Customer"

type EmailAddressConfirmed struct {
	id           *values.ID
	emailAddress *values.EmailAddress

	meta *shared.DomainEventMeta
}

/*** Factory Methods ***/

func EmailAddressWasConfirmed(
	id *values.ID,
	emailAddress *values.EmailAddress,
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

func (emailAddressConfirmed *EmailAddressConfirmed) ID() *values.ID {
	return emailAddressConfirmed.id
}

func (emailAddressConfirmed *EmailAddressConfirmed) EmailAddress() *values.EmailAddress {
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
		ID           *values.ID              `json:"id"`
		EmailAddress *values.EmailAddress    `json:"emailAddress"`
		Meta         *shared.DomainEventMeta `json:"meta"`
	}{
		ID:           emailAddressConfirmed.id,
		EmailAddress: emailAddressConfirmed.emailAddress,
		Meta:         emailAddressConfirmed.meta,
	}

	return json.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (emailAddressConfirmed *EmailAddressConfirmed) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		ID           *values.ID              `json:"id"`
		EmailAddress *values.EmailAddress    `json:"emailAddress"`
		Meta         *shared.DomainEventMeta `json:"meta"`
	}{}

	if err := json.Unmarshal(data, unmarshaledData); err != nil {
		return err
	}

	emailAddressConfirmed.id = unmarshaledData.ID
	emailAddressConfirmed.emailAddress = unmarshaledData.EmailAddress
	emailAddressConfirmed.meta = unmarshaledData.Meta

	return nil
}
