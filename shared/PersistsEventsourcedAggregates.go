package shared

type PersistsEventsourcedAggregates interface {
	Persist(aggregate EventsourcedAggregate) error
}
