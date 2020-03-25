package values

import (
	"regexp"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

var (
	emailAddressRegExp = regexp.MustCompile(`^[^\s]+@[^\s]+\.[\w]{2,}$`)
)

type EmailAddress struct {
	value string
}

func BuildEmailAddress(input string) (EmailAddress, error) {
	if matched := emailAddressRegExp.MatchString(input); !matched {
		err := lib.MarkAndWrapError(
			errors.New("input has invalid format"),
			lib.ErrInputIsInvalid,
			"BuildEmailAddress",
		)

		return EmailAddress{}, err
	}

	emailAddress := EmailAddress{value: input}

	return emailAddress, nil
}

func RebuildEmailAddress(input string) EmailAddress {
	return EmailAddress{value: input}
}

func (emailAddress EmailAddress) EmailAddress() string {
	return emailAddress.value
}

func (emailAddress EmailAddress) Equals(other EmailAddress) bool {
	return emailAddress.value == other.value
}