package es

type EventStore interface {
	AppendEventsToStream(streamID StreamID, events DomainEvents) error
	LoadEventStream(streamID StreamID, fromVersion uint, maxEvents uint) (DomainEvents, error)
	PurgeEventStream(streamID StreamID) error
}
