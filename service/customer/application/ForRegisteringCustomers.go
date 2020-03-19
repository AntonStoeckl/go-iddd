package application

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForRegisteringCustomers func(command commands.RegisterCustomer) error
