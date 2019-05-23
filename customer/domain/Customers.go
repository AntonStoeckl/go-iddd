package domain

import "go-iddd/customer/domain/values"

type Customers interface {
	Save(Customer) error
	FindBy(id *values.ID) (Customer, error)
}
