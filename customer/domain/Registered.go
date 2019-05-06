package domain

import (
	"encoding/json"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

const registeredAggregateName = "Customer"

type Registered interface {
	ID() *valueobjects.ID
	ConfirmableEmailAddress() *valueobjects.ConfirmableEmailAddress
	PersonName() *valueobjects.PersonName

	shared.DomainEvent
}

type registered struct {
	id                      *valueobjects.ID
	confirmableEmailAddress *valueobjects.ConfirmableEmailAddress
	personName              *valueobjects.PersonName

	meta *shared.DomainEventMeta
}

func ItWasRegistered(
	id *valueobjects.ID,
	confirmableEmailAddress *valueobjects.ConfirmableEmailAddress,
	personName *valueobjects.PersonName,
) *registered {

	registered := &registered{
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

func (registered *registered) ID() *valueobjects.ID {
	return registered.id
}

func (registered *registered) ConfirmableEmailAddress() *valueobjects.ConfirmableEmailAddress {
	return registered.confirmableEmailAddress
}

func (registered *registered) PersonName() *valueobjects.PersonName {
	return registered.personName
}

func (registered *registered) Identifier() string {
	return registered.meta.Identifier
}

func (registered *registered) EventName() string {
	return registered.meta.EventName
}

func (registered *registered) OccurredAt() string {
	return registered.meta.OccurredAt
}

func (registered *registered) MarshalJSON() ([]byte, error) {
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

func UnmarshalRegisteredFromJSON(jsonData []byte) (Registered, error) {
	var err error
	var data map[string]interface{}

	var id *valueobjects.ID
	var confirmableEmailAddress *valueobjects.ConfirmableEmailAddress
	var personName *valueobjects.PersonName
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
		case "confirmableEmailAddress":
			if confirmableEmailAddress, err = valueobjects.UnmarshalConfirmableEmailAddress(value); err != nil {
				return nil, err
			}
		case "personName":
			if personName, err = valueobjects.UnmarshalPersonName(value); err != nil {
				return nil, err
			}
		case "meta":
			if meta, err = shared.UnmarshalDomainEventMeta(value); err != nil {
				return nil, err
			}
		}
	}

	registered := &registered{
		id:                      id,
		confirmableEmailAddress: confirmableEmailAddress,
		personName:              personName,
		meta:                    meta,
	}

	return registered, nil
}
