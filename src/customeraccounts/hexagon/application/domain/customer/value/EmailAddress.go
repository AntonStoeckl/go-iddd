package value

type EmailAddress interface {
	String() string
	Equals(other EmailAddress) bool
}
