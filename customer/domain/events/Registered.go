package events

import (
	"encoding/json"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

const registeredAggregateName = "Customer"

type Registered struct {
	id                      *values.ID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName

	meta *shared.DomainEventMeta
}

/*** Factory Methods ***/

func ItWasRegistered(
	id *values.ID,
	confirmableEmailAddress *values.ConfirmableEmailAddress,
	personName *values.PersonName,
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

/*** Getter Methods ***/

func (registered *Registered) ID() *values.ID {
	return registered.id
}

func (registered *Registered) ConfirmableEmailAddress() *values.ConfirmableEmailAddress {
	return registered.confirmableEmailAddress
}

func (registered *Registered) PersonName() *values.PersonName {
	return registered.personName
}

/*** Implement shared.DomainEvent ***/

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
		ID                      *values.ID                      `json:"id"`
		ConfirmableEmailAddress *values.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		PersonName              *values.PersonName              `json:"personName"`
		Meta                    *shared.DomainEventMeta         `json:"meta"`
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
	unmarshaledData := &struct {
		ID                      *values.ID                      `json:"id"`
		ConfirmableEmailAddress *values.ConfirmableEmailAddress `json:"confirmableEmailAddress"`
		PersonName              *values.PersonName              `json:"personName"`
		Meta                    *shared.DomainEventMeta         `json:"meta"`
	}{}

	if err := json.Unmarshal(data, unmarshaledData); err != nil {
		return xerrors.Errorf("registered.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	registered.id = unmarshaledData.ID
	registered.confirmableEmailAddress = unmarshaledData.ConfirmableEmailAddress
	registered.personName = unmarshaledData.PersonName
	registered.meta = unmarshaledData.Meta

	return nil
}
