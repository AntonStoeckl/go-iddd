package lib

type PurgesEventStreams interface {
	PurgeEventStream(streamID StreamID) error
}
