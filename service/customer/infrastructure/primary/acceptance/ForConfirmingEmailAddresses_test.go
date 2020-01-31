package acceptance_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_ConfirmEmailAddress(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer, err := infrastructure.SetUpDIContainer()
		So(err, ShouldBeNil)
		commandHandler := diContainer.GetCustomerCommandHandler()

		Convey("Given a registered Customer", func() {
			register, err := commands.NewRegister(
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)
			So(err, ShouldBeNil)

			Convey("When a Customer's emailAddress is confirmed with a valid confirmationHash", func() {
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					register.CustomerID().ID(),
					register.EmailAddress().EmailAddress(),
					register.ConfirmationHash().Hash(),
				)
				So(err, ShouldBeNil)

				err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)

					Convey("And when this emailAddress is confirmed again", func() {
						confirmEmailAddress, err := commands.NewConfirmEmailAddress(
							register.CustomerID().ID(),
							register.EmailAddress().EmailAddress(),
							register.ConfirmationHash().Hash(),
						)
						So(err, ShouldBeNil)

						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

						Convey("It should succeed", func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})

			Convey("When a Customer's emailAddress is confirmed with an invalid confirmationHash", func() {
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					register.CustomerID().ID(),
					register.EmailAddress().EmailAddress(),
					"some_invalid_hash",
				)
				So(err, ShouldBeNil)

				err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)
				})
			})

			Convey("When a Customer's emailAddress is confirmed with an invalid command", func() {
				err := commandHandler.ConfirmEmailAddress(commands.ConfirmEmailAddress{})

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})

			err = diContainer.GetCustomerEventStore().Delete(register.CustomerID())
			So(err, ShouldBeNil)
		})

		Convey("Given an unregistered Customer", func() {
			Convey("When a Customer's emailAddress is confirmed", func() {
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					values.GenerateCustomerID().ID(),
					"john@doe.com",
					"any_hash",
				)
				So(err, ShouldBeNil)

				err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_ConfirmEmailAddress_WithErrorFromCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		customers := new(mocked.ForStoringCustomerEvents)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		commandHandler := application.NewCommandHandler(customers, db)

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")

		eventStream := es.DomainEvents{
			events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1)),
		}

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			emailAddress.EmailAddress(),
			confirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		Convey("Given the Customer can't be persisted because of a technical error", func() {
			customers.On("EventStreamFor", mock.Anything).Return(eventStream, nil).Once()
			customers.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(lib.ErrTechnical).Once()

			Convey("When a Customer's emailAddress is confirmed", func() {
				err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})
	})
}
