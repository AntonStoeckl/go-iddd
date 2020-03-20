package application

import (
	"go-iddd/service/customer/domain/customer/commands"
)

type ForConfirmingCustomerEmailAddresses func(command commands.ConfirmCustomerEmailAddress) error
