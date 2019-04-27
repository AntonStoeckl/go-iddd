package domain

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

type EmailAddressConfirmed interface {
	ID() valueobjects.ID
	EmailAddress() valueobjects.EmailAddress
	ConfirmationHash() valueobjects.ConfirmationHash

	shared.DomainEvent
}

type emailAddressConfirmed struct {
	id               valueobjects.ID
	emailAddress     valueobjects.EmailAddress
	confirmationHash valueobjects.ConfirmationHash

	meta *shared.DomainEventMeta
}

func EmailAddressWasConfirmed(
	id valueobjects.ID,
	emailAddress valueobjects.EmailAddress,
	confirmationHash valueobjects.ConfirmationHash,
) *emailAddressConfirmed {

	emailAddressConfirmed := &emailAddressConfirmed{
		id:               id,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	emailAddressConfirmed.meta = shared.NewDomainEventMeta(id, NewUnregisteredCustomer(), emailAddressConfirmed)

	return emailAddressConfirmed
}

func (emailAddressConfirmed *emailAddressConfirmed) ID() valueobjects.ID {
	return emailAddressConfirmed.id
}

func (emailAddressConfirmed *emailAddressConfirmed) EmailAddress() valueobjects.EmailAddress {
	return emailAddressConfirmed.emailAddress
}

func (emailAddressConfirmed *emailAddressConfirmed) ConfirmationHash() valueobjects.ConfirmationHash {
	return emailAddressConfirmed.confirmationHash
}

func (emailAddressConfirmed *emailAddressConfirmed) Identifier() shared.AggregateIdentifier {
	return emailAddressConfirmed.meta.Identifier
}

func (emailAddressConfirmed *emailAddressConfirmed) EventName() string {
	return emailAddressConfirmed.meta.EventName
}

func (emailAddressConfirmed *emailAddressConfirmed) OccurredAt() string {
	return emailAddressConfirmed.meta.OccurredAt
}
