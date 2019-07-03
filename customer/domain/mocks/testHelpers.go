package mocks

import (
	"go-iddd/shared"
	"reflect"
)

func FindCustomerEventIn(recordedEvents shared.DomainEvents, expectedEvent shared.DomainEvent) shared.DomainEvent {
	for _, event := range recordedEvents {
		if reflect.TypeOf(event) == reflect.TypeOf(expectedEvent) {
			return event
		}
	}

	return nil
}
