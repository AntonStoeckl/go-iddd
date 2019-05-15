package commands_test

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestNewConfirmEmailAddress(t *testing.T) {
	Convey("Given valid ID, EmailAddress and ConfirmationHash", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())

		Convey("When a new ConfirmEmailAddress command is created", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(id, emailAddress, confirmationHash)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(confirmEmailAddress, ShouldHaveSameTypeAs, (*commands.ConfirmEmailAddress)(nil))
			})
		})

		Convey("Given that ID is nil instead", func() {
			var id *values.ID
			conveyNewConfirmEmailAddressWithInvalidInput(id, emailAddress, confirmationHash)
		})

		Convey("Given that EmailAddress is nil instead", func() {
			var emailAddress *values.EmailAddress
			conveyNewConfirmEmailAddressWithInvalidInput(id, emailAddress, confirmationHash)
		})

		Convey("Given that ConfirmationHash is nil instead", func() {
			var confirmationHash *values.ConfirmationHash
			conveyNewConfirmEmailAddressWithInvalidInput(id, emailAddress, confirmationHash)
		})
	})
}

func conveyNewConfirmEmailAddressWithInvalidInput(
	id *values.ID,
	emailAddress *values.EmailAddress,
	confirmationHash *values.ConfirmationHash,
) {

	Convey("When a new ConfirmEmailAddress command is created", func() {
		confirmEmailAddress, err := commands.NewConfirmEmailAddress(id, emailAddress, confirmationHash)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(xerrors.Is(err, shared.ErrNilInput), ShouldBeTrue)
			So(confirmEmailAddress, ShouldBeNil)
		})
	})
}

func TestConfirmEmailAddressExposesExpectedValues(t *testing.T) {
	Convey("Given a ConfirmEmailAddress command", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(id, emailAddress, confirmationHash)
		So(err, ShouldBeNil)

		Convey("It should expose the expected values", func() {
			So(confirmEmailAddress.ID(), ShouldResemble, id)
			So(confirmEmailAddress.EmailAddress(), ShouldResemble, emailAddress)
			So(confirmEmailAddress.ConfirmationHash(), ShouldResemble, confirmationHash)
			So(confirmEmailAddress.CommandName(), ShouldEqual, "ConfirmEmailAddress")
			So(confirmEmailAddress.AggregateIdentifier(), ShouldResemble, id)
		})
	})
}
