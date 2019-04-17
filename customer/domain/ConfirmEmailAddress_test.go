package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/valueobjects"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfirmEmailAddress(t *testing.T) {
	Convey("Given that ID, EmailAddress and PersonName are valid", t, func() {
		id := valueobjects.GenerateID()
		emailAddress := valueobjects.ReconstituteEmailAddress("foo@bar.com")
		confirmationHash := valueobjects.GenerateConfirmationHash(emailAddress.String())

		Convey("When NewConfirmEmailAddress is invoked", func() {
			confirmEmailAddress, err := domain.NewConfirmEmailAddress(id, emailAddress, confirmationHash)

			Convey("Then it should create a ConfirmEmailAddress command", func() {
				So(err, ShouldBeNil)
				So(confirmEmailAddress, ShouldImplement, (*domain.ConfirmEmailAddress)(nil))
			})

			Convey("And then it should expose the expected CommandName, Identifier, ID, EmailAddress and ConfirmationHash ", func() {
				So(confirmEmailAddress.CommandName(), ShouldEqual, "ConfirmEmailAddress")
				So(confirmEmailAddress.Identifier(), ShouldEqual, id.String())
				So(confirmEmailAddress.ID(), ShouldEqual, id)
				So(confirmEmailAddress.EmailAddress(), ShouldEqual, emailAddress)
				So(confirmEmailAddress.ConfirmationHash(), ShouldEqual, confirmationHash)
			})
		})
	})

	Convey("Given that ID is nil", t, func() {
		var id valueobjects.ID
		emailAddress := valueobjects.ReconstituteEmailAddress("foo@bar.com")
		confirmationHash := valueobjects.GenerateConfirmationHash(emailAddress.String())

		conveyNewConfirmEmailAddressWithInvalidInput(id, emailAddress, confirmationHash)
	})

	Convey("Given that EmailAddress is nil", t, func() {
		id := valueobjects.GenerateID()
		var emailAddress valueobjects.EmailAddress
		confirmationHash := valueobjects.GenerateConfirmationHash("foo@bar.com")

		conveyNewConfirmEmailAddressWithInvalidInput(id, emailAddress, confirmationHash)
	})

	Convey("Given that PersonName is nil", t, func() {
		id := valueobjects.GenerateID()
		emailAddress := valueobjects.ReconstituteEmailAddress("foo@bar.com")
		var confirmationHash valueobjects.ConfirmationHash

		conveyNewConfirmEmailAddressWithInvalidInput(id, emailAddress, confirmationHash)
	})
}

func conveyNewConfirmEmailAddressWithInvalidInput(
	id valueobjects.ID,
	emailAddress valueobjects.EmailAddress,
	confirmationHash valueobjects.ConfirmationHash,
) {

	Convey("When NewConfirmEmailAddress is invoked", func() {
		confirmEmailAddress, err := domain.NewConfirmEmailAddress(id, emailAddress, confirmationHash)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(confirmEmailAddress, ShouldBeNil)
		})
	})
}
