package shared

type Aggregate interface {
	AggregateID() AggregateID
	AggregateName() string
}
