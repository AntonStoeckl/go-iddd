package shared

type Command interface {
	AggregateIdentifier() AggregateIdentifier
	CommandName() string
}
