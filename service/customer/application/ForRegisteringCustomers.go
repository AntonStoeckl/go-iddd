package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
)

type ForRegisteringCustomers func(command commands.RegisterCustomer) error
