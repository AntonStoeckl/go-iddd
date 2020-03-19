package application

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForConfirmingCustomerEmailAddresses func(command commands.ConfirmCustomerEmailAddress) error
