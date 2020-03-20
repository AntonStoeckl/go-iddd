package application

import (
	"go-iddd/service/customer/domain/customer/commands"
)

type ForRegisteringCustomers func(command commands.RegisterCustomer) error
