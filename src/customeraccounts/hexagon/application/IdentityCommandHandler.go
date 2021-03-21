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
	db                   *sql.DB
	uniqueIdentities     ForStoringUniqueIdentitiesWithTx
	identityEventStreams ForStoringIdentityEventStreamsWithTx
	maxRetries           uint8
}

func NewIdentityCommandHandler(
	db *sql.DB,
	uniqueIdentities ForStoringUniqueIdentitiesWithTx,
	identityEventStreams ForStoringIdentityEventStreamsWithTx,
	maxRetriesOnConcurrencyConflict uint8,
) *IdentityCommandHandler {

	return &IdentityCommandHandler{
		db:                   db,
		uniqueIdentities:     uniqueIdentities,
		identityEventStreams: identityEventStreams,
		maxRetries:           maxRetriesOnConcurrencyConflict,
	}
}

func (h *IdentityCommandHandler) HandleRegisterIdentity(
	identityIDValue value.IdentityID,
	emailAddress string,
	plainPassword string,
) error {

	registerIdentityWithSessions := func(
		uniqueIdentities ForStoringUniqueIdentities,
		identityEventStreams ForStoringIdentityEventStreams,
	) error {

		command, err := domain.BuildRegisterIdentity(
			identityIDValue,
			emailAddress,
			plainPassword,
		)
		if err != nil {
			return err
		}

		if err := uniqueIdentities.AddIdentity(identityIDValue, command.EmailAddress()); err != nil {
			return err
		}

		identityRegistered := identity.Register(command)

		if err := identityEventStreams.StartEventStream(identityRegistered); err != nil {
			return err
		}

		return nil
	}

	if err := h.handle(registerIdentityWithSessions); err != nil {
		return errors.Wrap(err, "IdentityCommandHandler.HandleRegisterIdentity")
	}

	return nil
}

type withSessions func(
	uniqueIdentities ForStoringUniqueIdentities,
	identityEventStreams ForStoringIdentityEventStreams,
) error

func (h *IdentityCommandHandler) handle(usecaseFn withSessions) error {
	return shared.RetryOnConcurrencyConflict(
		func() error {
			return shared.WrapInTx(
				func(tx *sql.Tx) error {
					return usecaseFn(
						h.uniqueIdentities.WithTx(tx),
						h.identityEventStreams.WithTx(tx),
					)
				},
				h.db,
			)
		},
		h.maxRetries,
	)
}

func (h *IdentityCommandHandler) ConfirmIdentityEmailAddress(
	identityID string,
	confirmationHash string,
) error {
	return errors.New("dummy error")
}
