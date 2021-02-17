package value

type ConfirmedEmailAddress string

func ToConfirmedEmailAddress(emailAddress EmailAddress) ConfirmedEmailAddress {
	return ConfirmedEmailAddress(emailAddress.String())
}

func RebuildConfirmedEmailAddress(input string) ConfirmedEmailAddress {
	return ConfirmedEmailAddress(input)
}

func (emailAddress ConfirmedEmailAddress) String() string {
	return string(emailAddress)
}

func (emailAddress ConfirmedEmailAddress) Equals(other EmailAddress) bool {
	return emailAddress.String() == other.String()
}
