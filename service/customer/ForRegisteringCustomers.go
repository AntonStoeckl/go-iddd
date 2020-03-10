package customer

import (
	"go-iddd/service/customer/application/domain/customer/commands"
)

type ForRegisteringCustomers func(command commands.RegisterCustomer) error
