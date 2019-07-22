package shared

type EventStore interface {
	AppendToStream(streamID *StreamID, events DomainEvents) error
	LoadEventStream(streamID *StreamID) (DomainEvents, error)
	LoadPartialEventStream(
		streamID *StreamID,
		fromVersion uint,
		maxEvents uint,
	) (DomainEvents, error)
}
