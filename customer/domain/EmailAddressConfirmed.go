package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

const emailAddressConfirmedAggregateName = "Customer"

type EmailAddressConfirmed interface {
	ID() *valueobjects.ID
	EmailAddress() *valueobjects.EmailAddress

	shared.DomainEvent
}

type emailAddressConfirmed struct {
	id           *valueobjects.ID
	emailAddress *valueobjects.EmailAddress

	meta *shared.DomainEventMeta
}

func EmailAddressWasConfirmed(
	id *valueobjects.ID,
	emailAddress *valueobjects.EmailAddress,
) *emailAddressConfirmed {

	emailAddressConfirmed := &emailAddressConfirmed{
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

func (emailAddressConfirmed *emailAddressConfirmed) ID() *valueobjects.ID {
	return emailAddressConfirmed.id
}

func (emailAddressConfirmed *emailAddressConfirmed) EmailAddress() *valueobjects.EmailAddress {
	return emailAddressConfirmed.emailAddress
}

func (emailAddressConfirmed *emailAddressConfirmed) Identifier() string {
	return emailAddressConfirmed.meta.Identifier
}

func (emailAddressConfirmed *emailAddressConfirmed) EventName() string {
	return emailAddressConfirmed.meta.EventName
}

func (emailAddressConfirmed *emailAddressConfirmed) OccurredAt() string {
	return emailAddressConfirmed.meta.OccurredAt
}

func (emailAddressConfirmed *emailAddressConfirmed) MarshalJSON() ([]byte, error) {
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

func UnmarshalEmailAddressConfirmedFromJSON(jsonData []byte) (EmailAddressConfirmed, error) {
	var err error
	var data map[string]interface{}

	var id *valueobjects.ID
	var emailAddress *valueobjects.EmailAddress
	var meta *shared.DomainEventMeta

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	for key, value := range data {
		switch key {
		case "id":
			if id, err = valueobjects.UnmarshalID(value); err != nil {
				return nil, err
			}
		case "emailAddress":
			if emailAddress, err = valueobjects.UnmarshalEmailAddress(value); err != nil {
				return nil, err
			}
		case "meta":
			if meta, err = shared.UnmarshalDomainEventMeta(value); err != nil {
				return nil, err
			}
		}
	}

	emailAddressConfirmed := &emailAddressConfirmed{
		id:           id,
		emailAddress: emailAddress,
		meta:         meta,
	}

	return emailAddressConfirmed, nil
}
