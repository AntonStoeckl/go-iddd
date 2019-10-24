package domain

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type Customers interface {
	Register(id *values.ID, recordedEvents shared.DomainEvents) error
	Of(id *values.ID) (*Customer, error)
}
