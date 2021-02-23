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
		emailAddress, err := value.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		plainPassword, err := value.BuildPlainPassword("superSecretPW123")
		So(err, ShouldBeNil)
		hashedPassword, err := value.HashedPasswordFromPlainPassword(plainPassword)
		So(err, ShouldBeNil)

		command := domain.BuildRegisterIdentity(
			identityID,
			emailAddress,
			hashedPassword,
		)

		Convey("\nSCENARIO: Register an Identity", func() {
			Convey("When RegisterIdentity", func() {
				event := identity.Register(command)

				Convey("Then IdentityRegistered", func() {
					So(event.IdentityID().Equals(identityID), ShouldBeTrue)
					So(event.EmailAddress().Equals(emailAddress), ShouldBeTrue)
					So(event.EmailAddress().ConfirmationHash().Equals(emailAddress.ConfirmationHash()), ShouldBeTrue)
					So(event.Password().Equals(hashedPassword), ShouldBeTrue)
					So(event.Password().CompareWith(plainPassword), ShouldBeTrue)
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
