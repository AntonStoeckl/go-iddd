package shared

type Aggregate interface {
	AggregateIdentifier() AggregateIdentifier
	AggregateName() string
}
