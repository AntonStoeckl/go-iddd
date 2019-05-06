package shared

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUnmarshaling = errors.New("unmarshaling failed")
)
