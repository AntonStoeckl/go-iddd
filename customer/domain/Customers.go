package domain

import "go-iddd/customer/domain/values"

//go:generate mockery -name Customers -output ../application/mocks -outpkg mocks -note "Regenerate by running `go generate` in customer/domain"

type Customers interface {
	Save(Customer) error
	FindBy(id *values.ID) (Customer, error)
}
