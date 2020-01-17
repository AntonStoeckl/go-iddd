package customer

import (
	"go-iddd/customer/domain/commands"
)

type ForRegisteringCustomers interface {
	Register(command commands.Register) error
}
