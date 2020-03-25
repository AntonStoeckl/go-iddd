package query_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerQueryHandler(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	queryHandler := diContainer.GetCustomerQueryHandler()

	Convey("\nSCENARIO: Trying to retrieve a non-existing Customer", t, func() {
		Convey("Given a Customer never registered", func() {
			Convey("When he tries to retrieve his account data", func() {
				customerID := values.GenerateCustomerID()
				view, err := queryHandler.CustomerViewByID(customerID)

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					So(view, ShouldBeZeroValue)
				})
			})
		})
	})
}
