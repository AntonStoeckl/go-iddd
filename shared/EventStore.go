package shared

type EventStore interface {
	LoadEventStream(streamID StreamID, fromVersion uint, maxEvents uint) (DomainEvents, error)
	AppendEventsToStream(streamID StreamID, events DomainEvents) error
}
