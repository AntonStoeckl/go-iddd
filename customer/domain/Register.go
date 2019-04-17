package domain

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

/*** The Customer behaviour method to apply the Register command ***/

func (customer *customer) register(register Register) {
	customer.id = register.ID()
	customer.emailAddress = register.ConfirmableEmailAddress()
	customer.personName = register.PersonName()
}

/*** The Register command itself - struct, factory, own getters, shared.Command getters ***/

type Register interface {
	ID() valueobjects.ID
	ConfirmableEmailAddress() valueobjects.ConfirmableEmailAddress
	PersonName() valueobjects.PersonName

	shared.Command
}

type register struct {
	id                      valueobjects.ID
	confirmableEmailAddress valueobjects.ConfirmableEmailAddress
	personName              valueobjects.PersonName
}

func NewRegister(
	id valueobjects.ID,
	emailAddress valueobjects.ConfirmableEmailAddress,
	personName valueobjects.PersonName,
) (*register, error) {

	command := &register{
		id:                      id,
		confirmableEmailAddress: emailAddress,
		personName:              personName,
	}

	if err := shared.AssertAllPropertiesAreNotNil(command); err != nil {
		return nil, err
	}

	return command, nil
}

func (register *register) ID() valueobjects.ID {
	return register.id
}

func (register *register) ConfirmableEmailAddress() valueobjects.ConfirmableEmailAddress {
	return register.confirmableEmailAddress
}

func (register *register) PersonName() valueobjects.PersonName {
	return register.personName
}

func (register *register) Identifier() string {
	return register.id.String()
}

func (register *register) CommandName() string {
	return shared.BuildCommandNameFor(register)
}
