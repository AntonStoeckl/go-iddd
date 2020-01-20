package lib

type Command interface {
	AggregateID() IdentifiesAggregates
	ShouldBeValid() error
}
