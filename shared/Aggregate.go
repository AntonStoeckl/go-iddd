package shared

type Aggregate interface {
	Execute(cmd Command) error
	AggregateID() IdentifiesAggregates
	AggregateName() string
}
