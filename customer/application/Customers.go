package application

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type Customers interface {
	Register(id values.CustomerID, recordedEvents shared.DomainEvents) error
	EventStream(id values.CustomerID) (shared.DomainEvents, error)
	Persist(id values.CustomerID, recordedEvents shared.DomainEvents) error
}
