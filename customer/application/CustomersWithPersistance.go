package application

import (
	"go-iddd/customer/domain"
	"go-iddd/shared"
)

type CustomersWithPersistance interface {
	domain.Customers
	shared.PersistsEventsourcedAggregates
}
