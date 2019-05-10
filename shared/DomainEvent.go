package shared

import (
	"fmt"
	"reflect"
)

type DomainEvent interface {
	Identifier() string
	EventName() string
	OccurredAt() string
}

func AssertEventPropertiesAreNotNilExcept(event DomainEvent, canBeNil ...string) error {
	elem := reflect.ValueOf(event).Elem()
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
