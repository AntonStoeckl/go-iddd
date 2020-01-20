package lib

type IdentifiesAggregates interface {
	String() string
	Equals(other IdentifiesAggregates) bool
}
