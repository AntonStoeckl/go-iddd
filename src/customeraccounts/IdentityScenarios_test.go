package customeraccounts_test

import (
	"fmt"
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/service/grpc"
	"github.com/AntonStoeckl/go-iddd/src/shared"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	. "github.com/smartystreets/goconvey/convey"
)

type identityScenarios struct {
	identityID          value.IdentityID
	otherIdentityID     value.IdentityID
	emailAddress        value.UnconfirmedEmailAddress
	changedEmailAddress value.UnconfirmedEmailAddress
	password            value.PlainPassword
	ea                  string // emailAddress
	cea                 string // changeEmailAddress
	ch                  string // confirmationHash
	cch                 string // changeConfirmationHash
	pw                  string // password

	registerIdentity            hexagon.ForRegisteringIdentities
	confirmIdentityEmailAddress hexagon.ForConfirmingIdentityEmailAddresses
	logIn                       hexagon.ForLoggingIn
}

func TestScenarios_ForRegisteringIdentities(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		v := initIdentityScenarios()

		Convey("\nSCENARIO: A prospective Customer registers his identity", func() {
			Convey(fmt.Sprintf("When a Customer registers his identity with [%s] and [%s]", v.ea, v.pw), func() {
				err := v.registerIdentity(v.identityID, v.ea, v.pw)
				So(err, ShouldBeNil)

				Convey(fmt.Sprintf("Then he should not be able to log in as his email address [%s] is unconfirmed", v.ea), func() {
					isLoggedIn, err := v.logIn(v.ea, v.pw)
					So(err, ShouldBeNil)
					So(isLoggedIn, ShouldBeFalse)
				})

				Convey(fmt.Sprintf("And when he confirms his email address [%s]", v.ea), func() {
					err := v.confirmIdentityEmailAddress(v.ea, v.ch)
					So(err, ShouldBeNil)

					Convey(fmt.Sprintf("Then he should be able to log in with [%s] and [%s]", v.ea, v.pw), func() {
						isLoggedIn, err := v.logIn(v.ea, v.pw)
						So(err, ShouldBeNil)
						So(isLoggedIn, ShouldBeTrue)
					})
				})
			})
		})
	})
}

func initIdentityScenarios() identityScenarios {
	customerID := value.GenerateIdentityID()
	otherCustomerID := value.GenerateIdentityID()
	emailAddress, err := value.BuildUnconfirmedEmailAddress("kevin@ball.net")
	So(err, ShouldBeNil)
	changedEmailAddress, err := value.BuildUnconfirmedEmailAddress("levinia@ball.net")
	So(err, ShouldBeNil)
	password, err := value.BuildPlainPassword(emailAddress.String())
	So(err, ShouldBeNil)

	logger := shared.NewNilLogger()
	config := grpc.MustBuildConfigFromEnv(logger)
	postgresDBConn := grpc.MustInitPostgresDB(config, logger)
	diContainer := grpc.MustBuildDIContainer(config, logger, grpc.UsePostgresDBConn(postgresDBConn))
	eventStore := diContainer.GetIdentityEventStore()

	identityCommandHandler := application.NewIdentityCommandHandler(eventStore)
	loginHandler := application.NewLoginHandler(eventStore)

	return identityScenarios{
		identityID:          customerID,
		otherIdentityID:     otherCustomerID,
		emailAddress:        emailAddress,
		changedEmailAddress: changedEmailAddress,
		password:            password,
		ea:                  emailAddress.String(),
		cea:                 changedEmailAddress.String(),
		ch:                  emailAddress.ConfirmationHash().String(),
		cch:                 changedEmailAddress.ConfirmationHash().String(),
		pw:                  password.String(),

		registerIdentity:            identityCommandHandler.RegisterIdentity,
		confirmIdentityEmailAddress: identityCommandHandler.ConfirmIdentityEmailAddress,
		logIn:                       loginHandler.Login,
	}
}
