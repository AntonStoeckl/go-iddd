package shared

type EventsourcedAggregate interface {
	StreamVersion() uint
	Apply(latestEvents DomainEvents)

	Aggregate
	RecordsEvents
}
