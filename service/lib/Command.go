package lib

// This is only a marker interface!
type Command interface {
	IsCommand() bool
}
