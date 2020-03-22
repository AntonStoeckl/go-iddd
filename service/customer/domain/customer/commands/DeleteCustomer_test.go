package commands_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildDeleteCustomerWithInvalidInput(t *testing.T) {
	Convey("When a DeleteCustomer command is built with an empty customerID", t, func() {
		_, err := commands.BuildCDeleteCustomer("")

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, lib.ErrInputIsInvalid), ShouldBeTrue)
		})
	})
}
