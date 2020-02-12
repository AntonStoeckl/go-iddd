package customer

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForRegistering func(command commands.Register) error
