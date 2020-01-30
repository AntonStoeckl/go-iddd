package es

type StreamID struct {
	value string
}

func NewStreamID(from string) StreamID {
	if from == "" {
		panic("newStreamID: empty input given")
	}

	return StreamID{value: from}
}

func (streamID StreamID) String() string {
	return streamID.value
}

func (streamID StreamID) Equals(other StreamID) bool {
	return streamID.value == other.value
}
