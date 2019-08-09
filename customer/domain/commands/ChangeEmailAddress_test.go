// Code generated by generate/main.go. DO NOT EDIT.

package commands_test

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewChangeEmailAddress(t *testing.T) {
	Convey("Given valid input", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"

		Convey("When a new ChangeEmailAddress command is created", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(id, emailAddress)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(changeEmailAddress, ShouldHaveSameTypeAs, (*commands.ChangeEmailAddress)(nil))
			})
		})

		Convey("Given that id is invalid", func() {
			id = ""
			conveyNewChangeEmailAddressWithInvalidInput(id, emailAddress)
		})

		Convey("Given that emailAddress is invalid", func() {
			emailAddress = ""
			conveyNewChangeEmailAddressWithInvalidInput(id, emailAddress)
		})
	})
}

func conveyNewChangeEmailAddressWithInvalidInput(
	id string,
	emailAddress string,
) {

	Convey("When a new ChangeEmailAddress command is created", func() {
		changeEmailAddress, err := commands.NewChangeEmailAddress(id, emailAddress)

		Convey("It should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrInputIsInvalid), ShouldBeTrue)
			So(changeEmailAddress, ShouldBeNil)
		})
	})
}

func TestChangeEmailAddressExposesExpectedValues(t *testing.T) {
	Convey("Given a ChangeEmailAddress command", t, func() {
		id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
		emailAddress := "john@doe.com"

		idValue, err := values.RebuildID(id)
		So(err, ShouldBeNil)
		emailAddressValue, err := values.NewEmailAddress(emailAddress)
		So(err, ShouldBeNil)

		changeEmailAddress, err := commands.NewChangeEmailAddress(id, emailAddress)
		So(err, ShouldBeNil)

		Convey("It should expose the expected values", func() {
			So(idValue.Equals(changeEmailAddress.ID()), ShouldBeTrue)
			So(emailAddressValue.Equals(changeEmailAddress.EmailAddress()), ShouldBeTrue)
			So(changeEmailAddress.CommandName(), ShouldEqual, "ChangeEmailAddress")
			So(idValue.Equals(changeEmailAddress.AggregateID()), ShouldBeTrue)
		})
	})
}
