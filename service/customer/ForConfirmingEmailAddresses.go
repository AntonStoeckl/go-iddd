package customer

import (
	"go-iddd/service/customer/domain/commands"
)

type ForConfirmingEmailAddresses interface {
	ConfirmEmailAddress(command commands.ConfirmEmailAddress) error
}
