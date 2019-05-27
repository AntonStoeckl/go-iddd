package shared

import "reflect"

type EventStream []DomainEvent

func (eventStream EventStream) FirstEventIsOfSameTypeAs(domainEvent DomainEvent) bool {
	if len(eventStream) == 0 {
		return false
	}

	return reflect.TypeOf(eventStream[0]) == reflect.TypeOf(domainEvent)
}
