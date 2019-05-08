package shared

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrNotEqual           = errors.New("not equal")
	ErrMarshalingFailed   = errors.New("marshaling failed")
	ErrUnmarshalingFailed = errors.New("unmarshaling failed")
)
