package application

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
)

type ForRegisteringCustomers func(command commands.RegisterCustomer) error
