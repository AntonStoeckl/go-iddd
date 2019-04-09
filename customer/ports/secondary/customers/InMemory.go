package customers

import (
	"go-iddd/customer/model"
	"go-iddd/customer/model/valueobjects"
)

type inMemory struct{}

func NewInMemoryCustomers() *inMemory {
	return &inMemory{}
}

func (customers *inMemory) Save(model.Customer) error {
	panic("implement me")
}

func (customers *inMemory) FindBy(id valueobjects.ID) (model.Customer, error) {
	panic("implement me")
}
