package shared

type IdentifiesAggregates interface {
	String() string
	Equals(other IdentifiesAggregates) bool
}
