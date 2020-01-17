package shared

type Command interface {
	AggregateID() IdentifiesAggregates
	ShouldBeValid() error
}
