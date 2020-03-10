package commands_test

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildChangeEmailAddressWithInvalidInput(t *testing.T) {
	Convey("When a ChangeCustomerEmailAddress command is built with an empty customerID", t, func() {
		_, err := commands.BuildChangeCustomerEmailAddress(
			"",
			"john@doe.com",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a ChangeCustomerEmailAddress command is built with an invalid emailAddress", t, func() {
		_, err := commands.BuildChangeCustomerEmailAddress(
			values.GenerateCustomerID().ID(),
			"foo@bar",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
