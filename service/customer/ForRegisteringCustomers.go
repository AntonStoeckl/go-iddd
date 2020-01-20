package customer

import (
	"go-iddd/service/customer/domain/commands"
)

type ForRegisteringCustomers interface {
	Register(command commands.Register) error
}
