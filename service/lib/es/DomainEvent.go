package es

type DomainEvent interface {
	EventName() string
	OccurredAt() string
	IndicatesAnError() bool
	StreamVersion() uint
}
