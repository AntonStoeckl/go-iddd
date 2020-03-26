package postgres_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAssertsUniqueEmailAddresses_With_Technical_Errors_From_DB(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	assertsUniqueEmailAddresses := diContainer.GetAssertsUniqueEmailAddresses()
	db := diContainer.GetPostgresDBConn()

	Convey("Setup", t, func() {
		tx, err := db.Begin()
		So(err, ShouldBeNil)

		var recordedEvents es.DomainEvents

		customerRegistered := events.CustomerWasRegistered(
			values.GenerateCustomerID(),
			values.RebuildEmailAddress("john@doe.com"),
			values.GenerateConfirmationHash("john@doe.com"),
			values.RebuildPersonName("John", "Doe"),
			1,
		)

		customerEmailAddressChanged := events.CustomerEmailAddressWasChanged(
			customerRegistered.CustomerID(),
			values.RebuildEmailAddress("john+changed@doe.com"),
			customerRegistered.ConfirmationHash(),
			customerRegistered.EmailAddress(),
			2,
		)

		customerDeleted := events.CustomerWasDeleted(
			customerRegistered.CustomerID(),
			customerEmailAddressChanged.EmailAddress(),
			3,
		)

		Convey("Given the DB transaction was already closed", func() {
			_ = tx.Rollback()

			Convey("When the uniqueness of the email address is asserted for a CustomerRegistered event", func() {
				recordedEvents = append(recordedEvents, customerRegistered)

				err := assertsUniqueEmailAddresses.Assert(recordedEvents, tx)

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When the uniqueness of the email address is asserted for a CustomerEmailAddressChanged event", func() {
				recordedEvents = append(recordedEvents, customerEmailAddressChanged)

				err := assertsUniqueEmailAddresses.Assert(recordedEvents, tx)

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When the uniqueness of the email address is asserted for a CustomerDeleted event", func() {
				recordedEvents = append(recordedEvents, customerDeleted)

				err := assertsUniqueEmailAddresses.Assert(recordedEvents, tx)

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When the unique email address of a Customer is removed", func() {
				err := assertsUniqueEmailAddresses.Remove(customerRegistered.CustomerID(), tx)

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}
