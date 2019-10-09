package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"sync"
)

type IdentityMap struct {
	customers sync.Map
}

func NewIdentityMap() *IdentityMap {
	return &IdentityMap{}
}

func (identityMap *IdentityMap) Memoize(customer domain.Customer) {
	identityMap.customers.Store(customer.AggregateID().String(), customer.Clone())
}

func (identityMap *IdentityMap) MemoizedCustomerOf(id *values.ID) (domain.Customer, bool) {
	customer, found := identityMap.customers.Load(id.String())
	if !found {
		return nil, false
	}

	return customer.(domain.Customer), true
}
