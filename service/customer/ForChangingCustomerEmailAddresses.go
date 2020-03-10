package customer

import (
	"go-iddd/service/customer/application/domain/customer/commands"
)

type ForChangingCustomerEmailAddresses func(command commands.ChangeCustomerEmailAddress) error
