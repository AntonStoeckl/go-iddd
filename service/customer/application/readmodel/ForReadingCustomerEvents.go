package readmodel

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/lib/es"
)

type ForReadingCustomerEvents interface {
	EventStreamFor(id values.CustomerID) (es.DomainEvents, error)
}
