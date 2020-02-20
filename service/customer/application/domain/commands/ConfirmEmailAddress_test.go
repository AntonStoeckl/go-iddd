package commands_test

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfirmEmailAddressWithInvalidInput(t *testing.T) {
	Convey("When a new ConfirmEmailAddress command is created with an empty customerID", t, func() {
		_, err := commands.NewConfirmEmailAddress(
			"",
			values.GenerateConfirmationHash("john@doe.com").Hash(),
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a new ConfirmEmailAddress command is created with an empty confirmationHash", t, func() {
		_, err := commands.NewConfirmEmailAddress(
			values.GenerateCustomerID().ID(),
			"",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
