package shared

type PersistsEventRecordingAggregates interface {
	Persist(aggregate EventRecordingAggregate) error
}
