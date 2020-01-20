package customer

import (
	"go-iddd/service/customer/application/domain/commands"
)

type ForConfirmingEmailAddresses interface {
	ConfirmEmailAddress(command commands.ConfirmEmailAddress) error
}
