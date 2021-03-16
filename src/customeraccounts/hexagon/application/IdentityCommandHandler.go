package application

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

type IdentityCommandHandler struct {
	db                              *sql.DB
	uniqueIdentities                ForStoringUniqueIdentitiesWithTx
	identityEventStreams            ForStoringIdentityEventStreamsWithTx
	maxRetriesOnConcurrencyConflict uint8
}

func NewIdentityCommandHandler(
	db *sql.DB,
	uniqueIdentities ForStoringUniqueIdentitiesWithTx,
	identityEventStreams ForStoringIdentityEventStreamsWithTx,
	maxRetriesOnConcurrencyConflict uint8,
) *IdentityCommandHandler {

	return &IdentityCommandHandler{
		db:                              db,
		uniqueIdentities:                uniqueIdentities,
		identityEventStreams:            identityEventStreams,
		maxRetriesOnConcurrencyConflict: maxRetriesOnConcurrencyConflict,
	}
}

func (h *IdentityCommandHandler) RegisterIdentity(
	identityIDValue value.IdentityID,
	emailAddress string,
	plainPassword string,
) error {

	wrapWithMsg := "IdentityCommandHandler.RegisterIdentity"

	emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	plainPasswordValue, err := value.BuildPlainPassword(plainPassword)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	hashedPasswordValue, err := value.HashedPasswordFromPlainPassword(plainPasswordValue)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	registerIdentity := domain.BuildRegisterIdentity(
		identityIDValue,
		emailAddressValue,
		hashedPasswordValue,
	)

	doRegister := func(tx *sql.Tx) error {
		uniqueIdentitiesSession := h.uniqueIdentities.WithTx(tx)
		identityEventStreamsSession := h.identityEventStreams.WithTx(tx)

		if err := uniqueIdentitiesSession.AddIdentity(identityIDValue, emailAddressValue); err != nil {
			return err
		}

		identityRegistered := identity.Register(registerIdentity)

		if err := identityEventStreamsSession.StartEventStream(identityRegistered); err != nil {
			return err
		}

		return nil
	}

	doRegisterWithinTx := func() error {
		return h.wrapWithTx(doRegister)
	}

	if err := shared.RetryOnConcurrencyConflict(doRegisterWithinTx, h.maxRetriesOnConcurrencyConflict); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *IdentityCommandHandler) wrapWithTx(fn func(tx *sql.Tx) error) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()

		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()

		return err
	}

	return nil
}

func (h *IdentityCommandHandler) ConfirmIdentityEmailAddress(
	customerID string,
	confirmationHash string,
) error {
	return errors.New("dummy error")
}
