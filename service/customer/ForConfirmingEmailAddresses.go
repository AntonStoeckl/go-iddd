package customer

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForConfirmingEmailAddresses func(command commands.ConfirmCustomerEmailAddress) error
