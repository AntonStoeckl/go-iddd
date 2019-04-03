package vo

type EmailAddress interface {
	EmailAddress() string
	IsConfirmed() bool
}

type emailAddress struct {
	value       string
	isConfirmed bool
}

func NewEmailAddress(value string) *emailAddress {
	newEmailAddress := &emailAddress{value: value}

	return newEmailAddress
}

func (emailAddress *emailAddress) EmailAddress() string {
	return emailAddress.value
}

func (emailAddress *emailAddress) IsConfirmed() bool {
	return emailAddress.isConfirmed
}
