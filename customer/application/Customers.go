package application

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type Customers interface {
	Register(id *values.CustomerID, recordedEvents shared.DomainEvents) error
	Of(id *values.CustomerID) (*domain.Customer, error)
	Persist(id *values.CustomerID, recordedEvents shared.DomainEvents) error
}
