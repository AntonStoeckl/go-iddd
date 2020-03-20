package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
)

type ForChangingCustomerEmailAddresses func(command commands.ChangeCustomerEmailAddress) error
