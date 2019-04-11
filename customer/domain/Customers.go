package domain

import "go-iddd/customer/domain/valueobjects"

type Customers interface {
	Save(Customer) error
	FindBy(id valueobjects.ID) (Customer, error)
}
