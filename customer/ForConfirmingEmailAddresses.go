package customer

import (
	"go-iddd/customer/domain/commands"
)

type ForConfirmingEmailAddresses interface {
	ConfirmEmailAddress(command commands.ConfirmEmailAddress) error
}
