package application

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
)

type Customers interface {
	Register(id values.CustomerID, recordedEvents lib.DomainEvents) error
	EventStream(id values.CustomerID) (lib.DomainEvents, error)
	Persist(id values.CustomerID, recordedEvents lib.DomainEvents) error
}
