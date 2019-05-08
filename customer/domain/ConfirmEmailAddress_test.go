package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfirmEmailAddress(t *testing.T) {
	Convey("Given valid ID, EmailAddress and ConfirmationHash", t, func() {
		id := values.GenerateID()
		emailAddress, err := values.NewEmailAddress("foo@bar.com")
		So(err, ShouldBeNil)
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())

		Convey("When a new ConfirmEmailAddress command is created", func() {
			confirmEmailAddress, err := domain.NewConfirmEmailAddress(id, emailAddress, confirmationHash)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(confirmEmailAddress, ShouldHaveSameTypeAs, (*domain.ConfirmEmailAddress)(nil))
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

		Convey("Given that PersonName is nil instead", func() {
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
		confirmEmailAddress, err := domain.NewConfirmEmailAddress(id, emailAddress, confirmationHash)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
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

		register, err := domain.NewConfirmEmailAddress(id, emailAddress, confirmationHash)
		So(err, ShouldBeNil)

		Convey("It should expose the expected values", func() {
			So(register.ID(), ShouldResemble, id)
			So(register.EmailAddress(), ShouldResemble, emailAddress)
			So(register.ConfirmationHash(), ShouldResemble, confirmationHash)
			So(register.CommandName(), ShouldEqual, "ConfirmEmailAddress")
			So(register.AggregateIdentifier(), ShouldResemble, id)
		})
	})
}
