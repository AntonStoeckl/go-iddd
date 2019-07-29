package shared

type PurgesEventStreams interface {
	PurgeEventStream(streamID *StreamID) error
}
