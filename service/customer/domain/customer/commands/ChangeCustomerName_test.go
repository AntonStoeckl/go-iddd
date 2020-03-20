package commands_test

import (
	"go-iddd/service/customer/domain/customer/commands"
	"go-iddd/service/customer/domain/customer/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildChangeNameWithInvalidInput(t *testing.T) {
	Convey("When a ChangeCustomerName command is built with an empty customerID", t, func() {
		_, err := commands.BuildChangeCustomerName(
			"",
			"John",
			"Doe",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a ChangeCustomerName command is built with an empty givenName", t, func() {
		_, err := commands.BuildChangeCustomerName(
			values.GenerateCustomerID().ID(),
			"",
			"Doe",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})

	Convey("When a ChangeCustomerName command is built with an empty familyName", t, func() {
		_, err := commands.BuildChangeCustomerName(
			values.GenerateCustomerID().ID(),
			"John",
			"",
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
