package domain

import "go-iddd/customer/domain/values"

type Customers interface {
	Register(customer *Customer) error
	Of(id *values.ID) (*Customer, error)
}
