package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
)

type ForChangingCustomerNames func(command commands.ChangeCustomerName) error
