package es

type DomainEvent interface {
	Meta() EventMeta
	IsFailureEvent() bool
	FailureReason() error
}
