package domain

import (
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

/*** Factory Method ***/

type Register struct {
	id           *valueobjects.ID
	emailAddress *valueobjects.EmailAddress
	personName   *valueobjects.PersonName
}

func NewRegister(
	id *valueobjects.ID,
	emailAddress *valueobjects.EmailAddress,
	personName *valueobjects.PersonName,
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

func (register *Register) ID() *valueobjects.ID {
	return register.id
}

func (register *Register) EmailAddress() *valueobjects.EmailAddress {
	return register.emailAddress
}

func (register *Register) PersonName() *valueobjects.PersonName {
	return register.personName
}

/*** Implement shared.Command ***/

func (register *Register) AggregateIdentifier() shared.AggregateIdentifier {
	return register.id
}

func (register *Register) CommandName() string {
	return shared.BuildCommandNameFor(register)
}
