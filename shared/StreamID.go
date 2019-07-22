package shared

import (
	"errors"

	"golang.org/x/xerrors"
)

type StreamID struct {
	value string
}

func NewStreamID(from string) (*StreamID, error) {
	newStreamID := &StreamID{value: from}

	if err := newStreamID.shouldBeValid(); err != nil {
		return nil, xerrors.Errorf("streamID.New: %s: %w", err, ErrInputIsInvalid)
	}

	return newStreamID, nil
}

func (streamID *StreamID) shouldBeValid() error {
	if streamID.value == "" {
		return errors.New("empty input for streamID")
	}

	return nil
}

func (streamID *StreamID) String() string {
	return streamID.value
}

func (streamID *StreamID) Equals(other *StreamID) bool {
	return streamID.value == other.value
}
