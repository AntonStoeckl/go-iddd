package application

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
)

type ForConfirmingCustomerEmailAddresses func(command commands.ConfirmCustomerEmailAddress) error
