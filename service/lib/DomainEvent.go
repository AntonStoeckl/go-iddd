package lib

type DomainEvent interface {
	EventName() string
	OccurredAt() string
	StreamVersion() uint
}
