package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"
)

type Register struct {
	id           *values.ID
	emailAddress *values.EmailAddress
	personName   *values.PersonName
}

/*** Factory Method ***/

func NewRegister(
	id string,
	emailAddress string,
	givenName string,
	familyName string,
) (*Register, error) {

	idValue, err := values.RebuildID(id)
	if err != nil {
		return nil, err
	}

	emailAddressValue, err := values.NewEmailAddress(emailAddress)
	if err != nil {
		return nil, err
	}

	personNameValue, err := values.NewPersonName(givenName, familyName)
	if err != nil {
		return nil, err
	}

	register := &Register{
		id:           idValue,
		emailAddress: emailAddressValue,
		personName:   personNameValue,
	}

	return register, nil
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
	commandType := reflect.TypeOf(register).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
