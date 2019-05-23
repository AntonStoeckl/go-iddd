package shared

import (
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
