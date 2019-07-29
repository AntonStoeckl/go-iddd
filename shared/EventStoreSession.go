package shared

type EventStoreSession interface {
	LoadEventStream(streamID *StreamID, fromVersion uint, maxEvents uint) (DomainEvents, error)
	AppendEventsToStream(streamID *StreamID, events DomainEvents) error
	Commit() error
	Rollback() error
}
