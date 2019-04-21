package domain

import (
    "go-iddd/customer/domain/valueobjects"
    "go-iddd/shared"
)

type Registered interface {
    ID() valueobjects.ID
    ConfirmableEmailAddress() valueobjects.ConfirmableEmailAddress
    PersonName() valueobjects.PersonName

    shared.DomainEvent
}

type registered struct {
    id                      valueobjects.ID
    confirmableEmailAddress valueobjects.ConfirmableEmailAddress
    personName              valueobjects.PersonName

    meta *shared.DomainEventMeta
}

func ItWasRegistered(
    id valueobjects.ID,
    confirmableEmailAddress valueobjects.ConfirmableEmailAddress,
    personName valueobjects.PersonName,
) *registered {

    registered := &registered{
        id:                      id,
        confirmableEmailAddress: confirmableEmailAddress,
        personName:              personName,
    }

    registered.meta = shared.NewDomainEventMeta(id, NewUnregisteredCustomer(), registered)

    return registered
}

func (registered *registered) ID() valueobjects.ID {
    return registered.id
}

func (registered *registered) ConfirmableEmailAddress() valueobjects.ConfirmableEmailAddress {
    return registered.confirmableEmailAddress
}

func (registered *registered) PersonName() valueobjects.PersonName {
    return registered.personName
}

func (registered *registered) Identifier() shared.AggregateIdentifier {
    return registered.meta.Identifier
}

func (registered *registered) EventName() string {
    return registered.meta.EventName
}

func (registered *registered) OccurredAt() string {
    return registered.meta.OccurredAt
}
