package identity

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
)

func Register(command domain.RegisterIdentity) domain.IdentityRegistered {
	event := domain.BuildIdentityRegistered(
		command.IdentityID(),
		command.EmailAddress(),
		command.MessageID(),
		1,
	)

	return event
}
