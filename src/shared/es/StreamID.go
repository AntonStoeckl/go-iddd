package es

type StreamID string

func BuildStreamID(from string) StreamID {
	if from == "" {
		panic("buildStreamID: empty input given")
	}

	return StreamID(from)
}

func (streamID StreamID) String() string {
	return string(streamID)
}
