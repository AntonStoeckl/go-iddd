package identity_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		identityID := value.GenerateIdentityID()
		emailAddress := "kevin@ball.com"
		plainPassword := "superSecretPW123"

		command, err := domain.BuildRegisterIdentity(
			identityID,
			emailAddress,
			plainPassword,
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO: Register an Identity", func() {
			Convey("When HandleRegisterIdentity", func() {
				event := identity.Register(command)

				Convey("Then IdentityRegistered", func() {
					So(event.IdentityID().Equals(identityID), ShouldBeTrue)
					So(event.EmailAddress().Equals(command.EmailAddress()), ShouldBeTrue)
					So(event.EmailAddress().ConfirmationHash().Equals(command.EmailAddress().ConfirmationHash()), ShouldBeTrue)
					So(event.Password().Equals(command.Password()), ShouldBeTrue)
					So(event.Password().CompareWith(value.RebuildPlainPassword(plainPassword)), ShouldBeTrue)
					So(event.IsFailureEvent(), ShouldBeFalse)
					So(event.FailureReason(), ShouldBeNil)
					So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
					So(event.Meta().MessageID(), ShouldNotBeEmpty)
					So(event.Meta().StreamVersion(), ShouldEqual, uint(1))
				})
			})
		})
	})
}
