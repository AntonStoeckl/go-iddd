package customers

import (
    "go-iddd/customer/domain"
    "go-iddd/customer/domain/valueobjects"
)

type inMemory struct{}

func NewInMemoryCustomers() *inMemory {
    return &inMemory{}
}

func (customers *inMemory) Save(domain.Customer) error {
    panic("implement me")
}

func (customers *inMemory) FindBy(id valueobjects.ID) (domain.Customer, error) {
    panic("implement me")
}
