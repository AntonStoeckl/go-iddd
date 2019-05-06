package shared

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrMarshaling   = errors.New("marshaling failed")
	ErrUnmarshaling = errors.New("unmarshaling failed")
)
