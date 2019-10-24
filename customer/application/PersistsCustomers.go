package application

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type PersistsCustomers interface {
	Persist(id *values.ID, recordedEvents shared.DomainEvents) error
}
