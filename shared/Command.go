package shared

type Command interface {
	AggregateID() AggregateID
	CommandName() string
}
