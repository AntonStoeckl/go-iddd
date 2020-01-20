package customer

import (
	"go-iddd/service/customer/domain/commands"
)

type ForChangingEmailAddresses interface {
	ChangeEmailAddress(command commands.ChangeEmailAddress) error
}
