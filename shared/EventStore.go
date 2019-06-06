package shared

type EventStore interface {
	AppendToStream(identifier AggregateID, events DomainEvents) error
	LoadEventStream(identifier AggregateID) (DomainEvents, error)
	LoadPartialEventStream(identifier AggregateID, fromVersion uint, maxEvents uint) (DomainEvents, error)
}
