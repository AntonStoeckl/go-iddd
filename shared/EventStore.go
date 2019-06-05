package shared

type EventStore interface {
	LoadEventStream(identifier AggregateID) (DomainEvents, error)
	AppendToStream(identifier AggregateID, events DomainEvents) error
}
