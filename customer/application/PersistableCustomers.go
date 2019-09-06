package application

import "go-iddd/customer/domain"

type PersistableCustomers interface {
	domain.Customers
	PersistsCustomers
}
