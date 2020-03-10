package commands_test

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildConfirmCustomerEmailAddressWithInvalidInput(t *testing.T) {
	Convey("When a ConfirmCustomerEmailAddress command is built with an empty customerID", t, func() {
		_, err := commands.BuildConfirmCustomerEmailAddress(
			"",
			values.GenerateConfirmationHash("john@doe.com").Hash(),
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a ConfirmCustomerEmailAddress command is built with an empty confirmationHash", t, func() {
		_, err := commands.BuildConfirmCustomerEmailAddress(
			values.GenerateCustomerID().ID(),
			"",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
