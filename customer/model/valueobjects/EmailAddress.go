package valueobjects

type EmailAddress interface {
	String() string
	Equals(other EmailAddress) bool
}

type emailAddress struct {
	value string
}

func NewEmailAddress(from string) *emailAddress {
	newEmailAddress := &emailAddress{
		value: from,
	}
	// TODO: validation

	return newEmailAddress
}

func (emailAddress *emailAddress) String() string {
	return emailAddress.value
}

func (emailAddress *emailAddress) Equals(other EmailAddress) bool {
	return emailAddress.String() == other.String()
}
