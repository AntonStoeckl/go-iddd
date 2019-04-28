package shared

import (
	"fmt"
	"reflect"
	"strings"
)

type Command interface {
	AggregateIdentifier() AggregateIdentifier
	CommandName() string
}

func BuildCommandNameFor(command Command) string {
	commandType := reflect.TypeOf(command).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}

func AssertAllPropertiesAreNotNil(command Command) error {
	return AssertPropertiesAreNotNilExcept(command)
}

func AssertPropertiesAreNotNilExcept(command Command, canBeNil ...string) error {
	elem := reflect.ValueOf(command).Elem()
	typeOf := elem.Type()

outer:
	for i := 0; i < elem.NumField(); i++ {
		property := elem.Field(i)
		propertyName := typeOf.Field(i).Name

		for _, bar := range canBeNil {
			if bar == propertyName {
				continue outer
			}
		}

		if property.IsNil() {
			return fmt.Errorf("nil given for: %s", propertyName)
		}
	}

	return nil
}
