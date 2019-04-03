package model

import (
	"go-iddd/person/model/command"
)

type Person interface {
	Register(register command.Register)
}

type person struct {
}

func NewPerson() *person {
	return &person{}
}

func (p *person) Register(register command.Register) {
	_ = register.ID()
	_ = register.EmailAddress()
	_ = register.Name()
}
