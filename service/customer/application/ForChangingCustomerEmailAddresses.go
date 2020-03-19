package application

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForChangingCustomerEmailAddresses func(command commands.ChangeCustomerEmailAddress) error
