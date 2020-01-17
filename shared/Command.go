package shared

type Command interface {
	AggregateID() IdentifiesAggregates
	CommandName() string
	ShouldBeValid() error
}
