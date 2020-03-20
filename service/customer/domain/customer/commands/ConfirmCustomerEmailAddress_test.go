package commands_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
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
