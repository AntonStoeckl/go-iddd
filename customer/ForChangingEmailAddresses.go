package customer

import (
	"go-iddd/customer/domain/commands"
)

type ForChangingEmailAddresses interface {
	ChangeEmailAddress(command commands.ChangeEmailAddress) error
}
