package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
)

type IdentityCommandHandler struct {
}

func NewIdentityCommandHandler() *IdentityCommandHandler {
	return &IdentityCommandHandler{}
}

func (h *IdentityCommandHandler) RegisterIdentity(
	identityIDValue value.IdentityID,
	emailAddress string,
	plainPassword string,
) error {
	return nil
}
