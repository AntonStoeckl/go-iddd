package customer

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForChangingEmailAddresses func(command commands.ChangeEmailAddress) error
