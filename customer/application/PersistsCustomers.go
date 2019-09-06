package application

import "go-iddd/customer/domain"

type PersistsCustomers interface {
	Persist(customer domain.Customer) error
}
