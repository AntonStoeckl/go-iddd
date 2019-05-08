package domain

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

/*** Factory Method ***/

type Register struct {
	id           *values.ID
	emailAddress *values.EmailAddress
	personName   *values.PersonName
}

func NewRegister(
	id *values.ID,
	emailAddress *values.EmailAddress,
	personName *values.PersonName,
) (*Register, error) {

	command := &Register{
		id:           id,
		emailAddress: emailAddress,
		personName:   personName,
	}

	if err := shared.AssertAllPropertiesAreNotNil(command); err != nil {
		return nil, err
	}

	return command, nil
}

/*** Getter Methods ***/

func (register *Register) ID() *values.ID {
	return register.id
}

func (register *Register) EmailAddress() *values.EmailAddress {
	return register.emailAddress
}

func (register *Register) PersonName() *values.PersonName {
	return register.personName
}

/*** Implement shared.Command ***/

func (register *Register) AggregateIdentifier() shared.AggregateIdentifier {
	return register.id
}

func (register *Register) CommandName() string {
	return shared.BuildCommandNameFor(register)
}
