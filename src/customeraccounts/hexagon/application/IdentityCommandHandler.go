package application

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/cockroachdb/errors"
)

type IdentityCommandHandler struct {
	db                   *sql.DB
	identities           ForStoringUniqueIdentities
	identityEventStreams ForStoringIdentityEventStreams
}

func NewIdentityCommandHandler(identityEventStreams ForStoringIdentityEventStreams) *IdentityCommandHandler {
	return &IdentityCommandHandler{
		identities:           nil,
		identityEventStreams: identityEventStreams,
	}
}

func (h *IdentityCommandHandler) RegisterIdentity(
	identityIDValue value.IdentityID,
	emailAddress string,
	plainPassword string,
) error {
	var err error
	var tx *sql.Tx

	tx, err = h.db.Begin()
	if err != nil {
		return err
	}

	foo, ok := h.identityEventStreams.(Transactional)
	if !ok {
		return errors.New("not transactional")
	}

	bar := foo.WithTx(tx)

	bar.RetrieveEventStream(identityIDValue)

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return errors.New("dummy error")
}

func (h *IdentityCommandHandler) ConfirmIdentityEmailAddress(
	customerID string,
	confirmationHash string,
) error {
	return errors.New("dummy error")
}
