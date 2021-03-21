package customeraccounts_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

type identityScenarios struct {
	// values as VO
	identityID          value.IdentityID
	otherIdentityID     value.IdentityID
	emailAddress        value.UnconfirmedEmailAddress
	changedEmailAddress value.UnconfirmedEmailAddress
	password            value.PlainPassword
	hashedPassword      value.HashedPassword

	// values as scalars
	ea  string // emailAddress
	cea string // changeEmailAddress
	ch  string // confirmationHash
	cch string // changeConfirmationHash
	pw  string // password
	hpw string // hashedPassword

	// usecases
	registerIdentity            hexagon.ForRegisteringIdentities
	confirmIdentityEmailAddress hexagon.ForConfirmingIdentityEmailAddresses
	logIn                       hexagon.ForLoggingIn

	// persistence
	db                 *sql.DB
	uniqueIdentities   application.ForStoringUniqueIdentitiesWithTx
	identityEventStore application.ForStoringIdentityEventStreamsWithTx
}

func TestScenarios_ForRegisteringIdentities(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		s := initIdentityScenarios()

		Convey("\nSCENARIO: A prospective Customer registers his identity", func() {
			Convey(fmt.Sprintf("When a Customer registers his identity with [%s] and [%s]", s.ea, s.pw), func() {
				err := s.registerIdentity(s.identityID, s.ea, s.pw)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then he should be able to log in with [%s] and [%s]", s.ea, s.pw), func() {
					isLoggedIn, err := s.logIn(s.ea, s.pw)
					So(err, ShouldBeNil)
					So(isLoggedIn, ShouldBeTrue)
				})
			})
		})

		Convey("\nSCENARIO: A prospective Customer can't register because his email address is already used", func() {
			Convey(fmt.Sprintf("Given a Customer registered his identity with [%s]", s.ea), func() {
				s.givenIdentityRegistered(s.identityID, s.emailAddress, s.hashedPassword)

				Convey(fmt.Sprintf("When another Customer registers with the same email address [%s]", s.ea), func() {
					err := s.registerIdentity(s.otherIdentityID, s.ea, s.pw)

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
					})
				})
			})
		})

		Reset(func() { s.reset() })
	})
}

func (s *identityScenarios) reset() {
	fn := func(tx *sql.Tx) error {
		uniqueIdentitiesSession := s.uniqueIdentities.WithTx(tx)
		eventStoreSession := s.identityEventStore.WithTx(tx)

		if err := uniqueIdentitiesSession.RemoveIdentity(s.identityID); err != nil {
			return err
		}

		if err := eventStoreSession.PurgeEventStream(s.identityID); err != nil {
			return err
		}

		if err := uniqueIdentitiesSession.RemoveIdentity(s.otherIdentityID); err != nil {
			return err
		}

		if err := eventStoreSession.PurgeEventStream(s.otherIdentityID); err != nil {
			return err
		}

		return nil
	}

	err := shared.WrapInTx(fn, s.db)
	So(err, ShouldBeNil)
}

func (s *identityScenarios) givenIdentityRegistered(
	identityID value.IdentityID,
	emailAddress value.UnconfirmedEmailAddress,
	password value.HashedPassword,
) {

	fn := func(tx *sql.Tx) error {
		event := domain.BuildIdentityRegistered(
			identityID,
			emailAddress,
			password,
			es.GenerateMessageID(),
			1,
		)
		uniqueIdentitiesSession := s.uniqueIdentities.WithTx(tx)
		eventStoreSession := s.identityEventStore.WithTx(tx)

		if err := uniqueIdentitiesSession.AddIdentity(identityID, emailAddress); err != nil {
			return err
		}

		if err := eventStoreSession.StartEventStream(event); err != nil {
			return err
		}

		return nil
	}

	err := shared.WrapInTx(fn, s.db)
	So(err, ShouldBeNil)
}

func initIdentityScenarios() *identityScenarios {
	logger := shared.NewNilLogger()
	config := grpc.MustBuildConfigFromEnv(logger)
	postgresDBConn := grpc.MustInitPostgresDB(config, logger)
	diContainer := grpc.MustBuildDIContainer(config, logger, grpc.UsePostgresDBConn(postgresDBConn))

	customerID := value.GenerateIdentityID()
	otherCustomerID := value.GenerateIdentityID()
	emailAddress, err := value.BuildUnconfirmedEmailAddress("kevin@ball.net")
	So(err, ShouldBeNil)
	changedEmailAddress, err := value.BuildUnconfirmedEmailAddress("levinia@ball.net")
	So(err, ShouldBeNil)
	password, err := value.BuildPlainPassword("jkjfU87aksdjf(&")
	So(err, ShouldBeNil)
	hashedPassword, err := value.HashedPasswordFromPlainPassword(password)
	So(err, ShouldBeNil)

	identityCommandHandler := diContainer.GetIdentityCommandHandler()
	loginHandler := diContainer.GetLoginHandler()

	return &identityScenarios{
		// values as VO
		identityID:          customerID,
		otherIdentityID:     otherCustomerID,
		emailAddress:        emailAddress,
		changedEmailAddress: changedEmailAddress,
		password:            password,
		hashedPassword:      hashedPassword,

		// values as scalars
		ea:  emailAddress.String(),
		cea: changedEmailAddress.String(),
		ch:  emailAddress.ConfirmationHash().String(),
		cch: changedEmailAddress.ConfirmationHash().String(),
		pw:  password.String(),
		hpw: hashedPassword.String(),

		// usecases
		registerIdentity:            identityCommandHandler.HandleRegisterIdentity,
		confirmIdentityEmailAddress: identityCommandHandler.ConfirmIdentityEmailAddress,
		logIn:                       loginHandler.HandleLogIn,

		// persistence
		db:                 diContainer.GetPostgresDBConn(),
		uniqueIdentities:   diContainer.GetUniqueIdentities(),
		identityEventStore: diContainer.GetIdentityEventStore(),
	}
}
