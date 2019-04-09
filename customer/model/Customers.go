package model

import "go-iddd/customer/model/valueobjects"

type Customers interface {
	Save(Customer) error
	FindBy(id valueobjects.ID) (Customer, error)
}
