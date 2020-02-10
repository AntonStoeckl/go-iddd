package commands_test

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewChangeEmailAddressWithInvalidInput(t *testing.T) {
	Convey("When a new ChangeEmailAddress command is created with an empty customerID", t, func() {
		_, err := commands.NewChangeEmailAddress(
			"",
			"john@doe.com",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a new ChangeEmailAddress command is created with an invalid emailAddress", t, func() {
		_, err := commands.NewChangeEmailAddress(
			values.GenerateCustomerID().ID(),
			"foo@bar",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
