package commands_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
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
