package application

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
)

type ForChangingCustomerEmailAddresses func(command commands.ChangeCustomerEmailAddress) error
