package query

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/mocked"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomerQueryHandler(t *testing.T) {
	customerEventStoreMock := new(mocked.ForStoringCustomerEvents)
	queryHandlerWithMock := NewCustomerQueryHandler(customerEventStoreMock)

	Convey("Prepare test artifacts", t, func() {
		customerID := values.GenerateCustomerID()

		Convey("\nSCENARIO: Technical problems with the CustomerEventStore", func() {
			Convey("Given a registered Customer", func() {
				Convey("and assuming the event stream can't be read", func() {
					customerEventStoreMock.
						On(
							"EventStreamFor",
							customerID,
						).
						Return(nil, lib.ErrTechnical).
						Once()

					Convey("When he tries to retrieve his account data", func() {
						view, err := queryHandlerWithMock.CustomerViewByID(customerID)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
							So(view, ShouldBeZeroValue)
						})
					})
				})
			})
		})
	})
}
