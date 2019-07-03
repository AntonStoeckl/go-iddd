package shared

type EventStore interface {
	AppendToStream(events DomainEvents) error
	LoadEventStream(identifier IdentifiesAggregates) (DomainEvents, error)
	LoadPartialEventStream(identifier IdentifiesAggregates, fromVersion uint, maxEvents uint) (DomainEvents, error)
}
