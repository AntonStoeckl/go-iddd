package application

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type Customers interface {
	Register(id *values.ID, recordedEvents shared.DomainEvents) error
	Of(id *values.ID) (*domain.Customer, error)
	Persist(id *values.ID, recordedEvents shared.DomainEvents) error
}
