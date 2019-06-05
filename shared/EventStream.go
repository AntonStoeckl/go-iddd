package shared

import (
	"errors"
	"fmt"
	"reflect"
)

type EventStream []DomainEvent

func (eventStream EventStream) FirstEventShouldBeOfSameTypeAs(domainEvent DomainEvent) error {
	if len(eventStream) == 0 {
		return errors.New("eventStream is empty")
	}

	expectedType := reflect.TypeOf(domainEvent)
	actualType := reflect.TypeOf(eventStream[0])

	if actualType != expectedType {
		return fmt.Errorf(
			"first event in eventStream should have type [%s] but has type [%s]",
			expectedType.String(),
			actualType.String(),
		)
	}

	return nil
}
