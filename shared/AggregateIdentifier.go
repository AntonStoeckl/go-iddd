package shared

type AggregateIdentifier interface {
    String() string
    Equals(other AggregateIdentifier) bool
}
