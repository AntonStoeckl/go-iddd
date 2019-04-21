package domain

import "go-iddd/customer/domain/valueobjects"

//go:generate mockery -name Customers -output ../application/mocks -outpkg mocks -note "Regenerate by running `go generate` in customer/domain"

type Customers interface {
    New() Customer
    Save(Customer) error
    FindBy(id valueobjects.ID) (Customer, error)
}
