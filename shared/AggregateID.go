package shared

type AggregateID interface {
	String() string
	Equals(other AggregateID) bool
}
