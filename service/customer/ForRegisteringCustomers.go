package customer

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForRegisteringCustomers interface {
	Register(command commands.Register) error
}
