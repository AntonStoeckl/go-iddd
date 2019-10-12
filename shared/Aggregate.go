package shared

type Aggregate interface {
	AggregateID() IdentifiesAggregates
	AggregateName() string
}
