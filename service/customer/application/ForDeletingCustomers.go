package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
)

type ForDeletingCustomers func(command commands.DeleteCustomer) error
