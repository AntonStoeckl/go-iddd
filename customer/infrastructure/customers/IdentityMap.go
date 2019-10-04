package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"sync"
)

type IdentityMap struct {
	customers map[string]domain.Customer
	mux       sync.Mutex
}

func NewIdentityMap() *IdentityMap {
	return &IdentityMap{
		customers: make(map[string]domain.Customer),
	}
}

func (identityMap *IdentityMap) Memoize(customer domain.Customer) {
	identityMap.mux.Lock()
	defer identityMap.mux.Unlock()

	identityMap.customers[customer.AggregateID().String()] = customer.Clone()
}

func (identityMap *IdentityMap) MemoizedCustomerOf(id *values.ID) (domain.Customer, bool) {
	identityMap.mux.Lock()
	defer identityMap.mux.Unlock()

	customer, found := identityMap.customers[id.String()]

	if !found {
		return nil, false
	}

	return customer, true
}
