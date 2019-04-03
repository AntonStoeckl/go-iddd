package persons

import (
	"go-iddd/person/model"
	"go-iddd/person/model/vo"
)

type inMemory struct{}

func NewInMemoryPersons() *inMemory {
	return &inMemory{}
}

func (persons *inMemory) Save(model.Person) error {
	panic("implement me")
}

func (persons *inMemory) GetBy(id vo.ID) (model.Person, error) {
	panic("implement me")
}
