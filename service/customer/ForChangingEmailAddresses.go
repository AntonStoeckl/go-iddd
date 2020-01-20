package customer

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForChangingEmailAddresses interface {
	ChangeEmailAddress(command commands.ChangeEmailAddress) error
}
