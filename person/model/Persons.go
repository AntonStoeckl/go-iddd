package model

import "go-iddd/person/model/vo"

type Persons interface {
	Save(Person) error
	GetBy(id vo.ID) (Person, error)
}
