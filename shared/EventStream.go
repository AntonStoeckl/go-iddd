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

	if reflect.TypeOf(eventStream[0]) != reflect.TypeOf(domainEvent) {
		return fmt.Errorf(
			"first event in eventStream is not of type [%s] but [%s]",
			domainEvent.EventName(),
			eventStream[0].EventName(),
		)
	}

	return nil
}
