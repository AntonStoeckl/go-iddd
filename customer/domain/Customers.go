package domain

import "go-iddd/customer/domain/values"

type Customers interface {
	Register(Customer) error
	Save(Customer) error
	Of(id *values.ID) (Customer, error)
}
