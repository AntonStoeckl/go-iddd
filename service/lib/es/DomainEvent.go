package es

type DomainEvent interface {
	EventName() string
	OccurredAt() string
	StreamVersion() uint
}
