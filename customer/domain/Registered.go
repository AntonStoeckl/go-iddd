package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

const registeredAggregateName = "Customer"

type Registered struct {
	id                      *valueobjects.ID
	confirmableEmailAddress *valueobjects.ConfirmableEmailAddress
	personName              *valueobjects.PersonName

	meta *shared.DomainEventMeta
}

func ItWasRegistered(
	id *valueobjects.ID,
	confirmableEmailAddress *valueobjects.ConfirmableEmailAddress,
	personName *valueobjects.PersonName,
) *Registered {

	registered := &Registered{
		id:                      id,
		confirmableEmailAddress: confirmableEmailAddress,
		personName:              personName,
	}

	registered.meta = shared.NewDomainEventMeta(
		id.String(),
		registered,
		registeredAggregateName,
	)

	return registered
}

func (registered *Registered) ID() *valueobjects.ID {
	return registered.id
}

func (registered *Registered) ConfirmableEmailAddress() *valueobjects.ConfirmableEmailAddress {
	return registered.confirmableEmailAddress
}

func (registered *Registered) PersonName() *valueobjects.PersonName {
	return registered.personName
}

func (registered *Registered) Identifier() string {
	return registered.meta.Identifier
}

func (registered *Registered) EventName() string {
	return registered.meta.EventName
}

func (registered *Registered) OccurredAt() string {
	return registered.meta.OccurredAt
}

/*** Implement json.Marshaler ***/

func (registered *Registered) MarshalJSON() ([]byte, error) {
	data := &struct {
		ID                      *valueobjects.ID                      `json:"id"`
		ConfirmableEmailAddress *valueobjects.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		PersonName              *valueobjects.PersonName              `json:"personName"`
		Meta                    *shared.DomainEventMeta               `json:"meta"`
	}{
		ID:                      registered.id,
		ConfirmableEmailAddress: registered.confirmableEmailAddress,
		PersonName:              registered.personName,
		Meta:                    registered.meta,
	}

	return json.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func (registered *Registered) UnmarshalJSON(data []byte) error {
	values := &struct {
		ID                      *valueobjects.ID                      `json:"id"`
		ConfirmableEmailAddress *valueobjects.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		PersonName              *valueobjects.PersonName              `json:"personName"`
		Meta                    *shared.DomainEventMeta               `json:"meta"`
	}{}

	if err := json.Unmarshal(data, values); err != nil {
		return err
	}

	registered.id = values.ID
	registered.confirmableEmailAddress = values.ConfirmableEmailAddress
	registered.personName = values.PersonName
	registered.meta = values.Meta

	return nil
}
