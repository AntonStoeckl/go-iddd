package command

import (
	"go-iddd/person/model/vo"
	"go-iddd/shared"
)

type Register interface {
	ID() vo.ID
	EmailAddress() vo.EmailAddress
	Name() vo.Name
}

type register struct {
	id           vo.ID
	emailAddress vo.EmailAddress
	name         vo.Name
}

func NewRegister(id vo.ID, emailAddress vo.EmailAddress, name vo.Name) (*register, error) {
	command := &register{
		id:           id,
		emailAddress: emailAddress,
		name:         name,
	}

	if err := shared.AssertAllPropertiesAreNotNil(command); err != nil {
		return nil, err
	}

	return command, nil
}

func (register *register) ID() vo.ID {
	return register.id
}

func (register *register) EmailAddress() vo.EmailAddress {
	return register.emailAddress
}

func (register *register) Name() vo.Name {
	return register.name
}

func (register *register) Identifier() string {
	return register.id.ID()
}

func (register *register) CommandName() string {
	return shared.BuildNameFor(register)
}
