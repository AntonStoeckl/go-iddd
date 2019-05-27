package shared

type DomainEvent interface {
	Identifier() string
	EventName() string
	OccurredAt() string
	StreamVersion() uint
}
