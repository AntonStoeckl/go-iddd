package commands_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildChangeCustomerEmailAddressWithInvalidInput(t *testing.T) {
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
