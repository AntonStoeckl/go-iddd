package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

type Register struct {
	id           *values.ID
	emailAddress *values.EmailAddress
	personName   *values.PersonName
}

/*** Factory Method ***/

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

	if err := shared.AssertAllCommandPropertiesAreNotNil(command); err != nil {
		return nil, xerrors.Errorf("register.New -> %s: %w", err, shared.ErrNilInput)
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
