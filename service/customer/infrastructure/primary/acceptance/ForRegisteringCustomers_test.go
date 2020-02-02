package acceptance_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_ForRegisteringCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)
		db, sqlMock, err := sqlmock.New()
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(customerEventStore, db)

		register, err := commands.NewRegister(
			"john@doe.com",
			"John",
			"Doe",
		)
		So(err, ShouldBeNil)

		containsOnlyRegisteredEvent := func(recordedEvents es.DomainEvents) bool {
			if len(recordedEvents) != 1 {
				return false
			}

			_, ok := recordedEvents[0].(events.Registered)

			return ok
		}

		Convey("Given no Customer with the same ID exists", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()

			customerEventStore.
				On(
					"CreateStreamFrom",
					mock.MatchedBy(containsOnlyRegisteredEvent),
					register.CustomerID(),
					mock.AnythingOfType("*sql.Tx"),
				).
				Return(nil).
				Once()

			Convey("When a Customer is registered", func() {
				err = commandHandler.Register(register)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})

		Convey("Given a Customer with the same ID already exists", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()

			customerEventStore.
				On(
					"CreateStreamFrom",
					mock.MatchedBy(containsOnlyRegisteredEvent),
					register.CustomerID(),
					mock.AnythingOfType("*sql.Tx"),
				).
				Return(lib.ErrDuplicate).
				Once()

			Convey("When a Customer is registered", func() {
				err = commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
					So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})

		Convey("When a Customer is registered with an invalid command", func() {
			err := commandHandler.Register(commands.Register{})

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
			})
		})

		Convey("Assuming that beginning the transaction fails", func() {
			sqlMock.ExpectBegin().WillReturnError(lib.ErrTechnical)

			Convey("When a Customer is registered", func() {
				err := commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})

		Convey("Assuming that committing the transaction fails", func() {
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit().WillReturnError(lib.ErrTechnical)

			customerEventStore.On("CreateStreamFrom", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

			Convey("When a Customer is registered", func() {
				err := commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})
	})
}
